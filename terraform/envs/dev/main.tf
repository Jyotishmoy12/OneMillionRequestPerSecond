data "aws_availability_zones" "available" {
  state = "available"
}

locals {
  name_prefix = "${var.project_name}-${var.environment}"
}

resource "aws_vpc" "main" {
  cidr_block           = var.vpc_cidr
  enable_dns_support   = true
  enable_dns_hostnames = true

  tags = {
    Name        = "${local.name_prefix}-vpc"
    Project     = var.project_name
    Environment = var.environment
  }
}

resource "aws_internet_gateway" "main" {
  vpc_id = aws_vpc.main.id

  tags = {
    Name        = "${local.name_prefix}-igw"
    Project     = var.project_name
    Environment = var.environment
  }
}

resource "aws_subnet" "public" {
  count = length(var.public_subnet_cidrs)

  vpc_id                  = aws_vpc.main.id
  cidr_block              = var.public_subnet_cidrs[count.index]
  availability_zone       = data.aws_availability_zones.available.names[count.index]
  map_public_ip_on_launch = true

  tags = {
    Name        = "${local.name_prefix}-public-${count.index + 1}"
    Project     = var.project_name
    Environment = var.environment
  }
}

resource "aws_route_table" "public" {
  vpc_id = aws_vpc.main.id

  tags = {
    Name        = "${local.name_prefix}-public-rt"
    Project     = var.project_name
    Environment = var.environment
  }
}

resource "aws_route" "public_internet" {
  route_table_id         = aws_route_table.public.id
  destination_cidr_block = "0.0.0.0/0"
  gateway_id             = aws_internet_gateway.main.id
}

resource "aws_route_table_association" "public" {
  count = length(aws_subnet.public)

  subnet_id      = aws_subnet.public[count.index].id
  route_table_id = aws_route_table.public.id
}

resource "aws_security_group" "load_balancer" {
  name        = "${local.name_prefix}-lb-sg"
  description = "Allow public HTTP traffic to the load balancer."
  vpc_id      = aws_vpc.main.id

  ingress {
    description = "HTTP from internet"
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    description = "All outbound traffic"
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name        = "${local.name_prefix}-lb-sg"
    Project     = var.project_name
    Environment = var.environment
  }
}

resource "aws_security_group" "api" {
  name        = "${local.name_prefix}-api-sg"
  description = "Allow traffic from the load balancer to the API."
  vpc_id      = aws_vpc.main.id

  ingress {
    description     = "API traffic from load balancer"
    from_port       = 8080
    to_port         = 8080
    protocol        = "tcp"
    security_groups = [aws_security_group.load_balancer.id]
  }

  egress {
    description = "All outbound traffic"
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name        = "${local.name_prefix}-api-sg"
    Project     = var.project_name
    Environment = var.environment
  }
}

resource "aws_security_group" "postgres" {
  name        = "${local.name_prefix}-postgres-sg"
  description = "Allow Postgres traffic from the API."
  vpc_id      = aws_vpc.main.id

  ingress {
    description     = "Postgres from API"
    from_port       = 5432
    to_port         = 5432
    protocol        = "tcp"
    security_groups = [aws_security_group.api.id]
  }

  egress {
    description = "All outbound traffic"
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name        = "${local.name_prefix}-postgres-sg"
    Project     = var.project_name
    Environment = var.environment
  }
}

resource "aws_security_group" "valkey" {
  name        = "${local.name_prefix}-valkey-sg"
  description = "Allow Valkey traffic from the API."
  vpc_id      = aws_vpc.main.id

  ingress {
    description     = "Valkey from API"
    from_port       = 6379
    to_port         = 6379
    protocol        = "tcp"
    security_groups = [aws_security_group.api.id]
  }

  egress {
    description = "All outbound traffic"
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name        = "${local.name_prefix}-valkey-sg"
    Project     = var.project_name
    Environment = var.environment
  }
}


resource "aws_security_group" "load_generator" {
  name        = "${local.name_prefix}-loadgen-sg"
  description = "Allow outbound traffic from the benchmark load generator."
  vpc_id      = aws_vpc.main.id

  egress {
    description = "All outbound traffic"
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name        = "${local.name_prefix}-loadgen-sg"
    Project     = var.project_name
    Environment = var.environment
  }
}


resource "aws_lb" "api" {
  name               = "${local.name_prefix}-api-alb"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [aws_security_group.load_balancer.id]
  subnets            = aws_subnet.public[*].id

  tags = {
    Name        = "${local.name_prefix}-api-alb"
    Project     = var.project_name
    Environment = var.environment
  }
}

resource "aws_lb_target_group" "api" {
  name        = "${local.name_prefix}-api-tg"
  port        = var.api_container_port
  protocol    = "HTTP"
  target_type = "instance"
  vpc_id      = aws_vpc.main.id

  health_check {
    enabled             = true
    path                = "/healthz"
    matcher             = "200"
    interval            = 15
    timeout             = 5
    healthy_threshold   = 2
    unhealthy_threshold = 2
  }

  tags = {
    Name        = "${local.name_prefix}-api-tg"
    Project     = var.project_name
    Environment = var.environment
  }
}

resource "aws_lb_listener" "api_http" {
  load_balancer_arn = aws_lb.api.arn
  port              = 80
  protocol          = "HTTP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.api.arn
  }
}


resource "aws_ecr_repository" "api" {
  name                 = "${local.name_prefix}-api"
  image_tag_mutability = "MUTABLE"
  force_delete         = true

  image_scanning_configuration {
    scan_on_push = true
  }

  encryption_configuration {
    encryption_type = "AES256"
  }

  tags = {
    Name        = "${local.name_prefix}-api"
    Project     = var.project_name
    Environment = var.environment
  }
}


resource "aws_ecs_cluster" "main" {
  name = "${local.name_prefix}-cluster"

  setting {
    name  = "containerInsights"
    value = "disabled"
  }

  tags = {
    Name        = "${local.name_prefix}-cluster"
    Project     = var.project_name
    Environment = var.environment
  }
}

resource "aws_cloudwatch_log_group" "api" {
  name              = "/ecs/${local.name_prefix}-api"
  retention_in_days = 7

  tags = {
    Name        = "${local.name_prefix}-api-logs"
    Project     = var.project_name
    Environment = var.environment
  }
}

resource "aws_iam_role" "ecs_task_execution" {
  name = "${local.name_prefix}-ecs-task-execution-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Service = "ecs-tasks.amazonaws.com"
        }
        Action = "sts:AssumeRole"
      }
    ]
  })

  tags = {
    Name        = "${local.name_prefix}-ecs-task-execution-role"
    Project     = var.project_name
    Environment = var.environment
  }
}

resource "aws_iam_role_policy_attachment" "ecs_task_execution" {
  role       = aws_iam_role.ecs_task_execution.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}


resource "aws_ecs_task_definition" "api" {
  family                   = "${local.name_prefix}-api"
  requires_compatibilities = ["EC2"]
  network_mode             = "host"
  cpu                      = tostring(var.api_cpu)
  memory                   = tostring(var.api_memory)
  execution_role_arn       = aws_iam_role.ecs_task_execution.arn

  container_definitions = jsonencode([
    {
      name      = "postgres"
      image     = "postgres:17-alpine"
      essential = true

      portMappings = [
        {
          containerPort = 5432
          hostPort      = 5432
          protocol      = "tcp"
        }
      ]

      environment = [
        {
          name  = "POSTGRES_USER"
          value = "postgres"
        },
        {
          name  = "POSTGRES_PASSWORD"
          value = "postgres"
        },
        {
          name  = "POSTGRES_DB"
          value = "onemillionrps"
        }
      ]

      memoryReservation = 128

      logConfiguration = {
        logDriver = "awslogs"
        options = {
          awslogs-group         = aws_cloudwatch_log_group.api.name
          awslogs-region        = var.aws_region
          awslogs-stream-prefix = "postgres"
        }
      }
    },
    {
      name      = "valkey"
      image     = "valkey/valkey:8-alpine"
      essential = true

      command = [
        "valkey-server",
        "--save",
        "",
        "--appendonly",
        "no"
      ]

      portMappings = [
        {
          containerPort = 6379
          hostPort      = 6379
          protocol      = "tcp"
        }
      ]

      memoryReservation = 64

      logConfiguration = {
        logDriver = "awslogs"
        options = {
          awslogs-group         = aws_cloudwatch_log_group.api.name
          awslogs-region        = var.aws_region
          awslogs-stream-prefix = "valkey"
        }
      }
    },
    {
      name      = "api"
      image     = "${aws_ecr_repository.api.repository_url}:${var.api_image_tag}"
      essential = true

      dependsOn = [
        {
          containerName = "postgres"
          condition     = "START"
        },
        {
          containerName = "valkey"
          condition     = "START"
        }
      ]

      portMappings = [
        {
          containerPort = var.api_container_port
          hostPort      = var.api_container_port
          protocol      = "tcp"
        }
      ]

      environment = [
        {
          name  = "PORT"
          value = tostring(var.api_container_port)
        },
        {
          name  = "DATABASE_URL"
          value = "postgres://postgres:postgres@localhost:5432/onemillionrps?sslmode=disable"
        },
        {
          name  = "REDIS_ADDR"
          value = "localhost:6379"
        },
        {
          name  = "DB_MAX_CONNS"
          value = "50"
        },
        {
          name  = "DB_MIN_CONNS"
          value = "10"
        }
      ]

      memoryReservation = 128

      logConfiguration = {
        logDriver = "awslogs"
        options = {
          awslogs-group         = aws_cloudwatch_log_group.api.name
          awslogs-region        = var.aws_region
          awslogs-stream-prefix = "api"
        }
      }
    }
  ])

  tags = {
    Name        = "${local.name_prefix}-api-task"
    Project     = var.project_name
    Environment = var.environment
  }
}

data "aws_ami" "ecs_optimized" {
  most_recent = true
  owners      = ["amazon"]

  filter {
    name   = "name"
    values = ["al2023-ami-ecs-hvm-*-x86_64"]
  }
}

data "aws_ami" "amazon_linux_2023" {
  most_recent = true
  owners      = ["amazon"]

  filter {
    name   = "name"
    values = ["al2023-ami-2023.*-x86_64"]
  }

  filter {
    name   = "state"
    values = ["available"]
  }
}

resource "aws_iam_role" "ecs_instance" {
  name = "${local.name_prefix}-ecs-instance-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Service = "ec2.amazonaws.com"
        }
        Action = "sts:AssumeRole"
      }
    ]
  })

  tags = {
    Name        = "${local.name_prefix}-ecs-instance-role"
    Project     = var.project_name
    Environment = var.environment
  }
}

resource "aws_iam_role_policy_attachment" "ecs_instance" {
  role       = aws_iam_role.ecs_instance.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonEC2ContainerServiceforEC2Role"
}

resource "aws_iam_instance_profile" "ecs_instance" {
  name = "${local.name_prefix}-ecs-instance-profile"
  role = aws_iam_role.ecs_instance.name
}

resource "aws_launch_template" "ecs" {
  name_prefix   = "${local.name_prefix}-ecs-"
  image_id      = data.aws_ami.ecs_optimized.id
  instance_type = var.ecs_instance_type

  iam_instance_profile {
    name = aws_iam_instance_profile.ecs_instance.name
  }

  vpc_security_group_ids = [aws_security_group.api.id]

  user_data = base64encode(<<-EOF
    #!/bin/bash
    echo ECS_CLUSTER=${aws_ecs_cluster.main.name} >> /etc/ecs/ecs.config
  EOF
  )

  tag_specifications {
    resource_type = "instance"

    tags = {
      Name        = "${local.name_prefix}-ecs-instance"
      Project     = var.project_name
      Environment = var.environment
    }
  }

  tag_specifications {
    resource_type = "volume"

    tags = {
      Name        = "${local.name_prefix}-ecs-volume"
      Project     = var.project_name
      Environment = var.environment
    }
  }

  tags = {
    Name        = "${local.name_prefix}-ecs-launch-template"
    Project     = var.project_name
    Environment = var.environment
  }
}

resource "aws_autoscaling_group" "ecs" {
  name                = "${local.name_prefix}-ecs-asg"
  min_size            = var.ecs_min_size
  max_size            = var.ecs_max_size
  desired_capacity    = var.ecs_desired_capacity
  vpc_zone_identifier = aws_subnet.public[*].id

  launch_template {
    id      = aws_launch_template.ecs.id
    version = "$Latest"
  }

  tag {
    key                 = "Name"
    value               = "${local.name_prefix}-ecs-asg"
    propagate_at_launch = false
  }

  tag {
    key                 = "AmazonECSManaged"
    value               = "true"
    propagate_at_launch = true
  }
}

resource "aws_ecs_capacity_provider" "ec2" {
  name = "${local.name_prefix}-ec2-capacity-provider"

  auto_scaling_group_provider {
    auto_scaling_group_arn         = aws_autoscaling_group.ecs.arn
    managed_termination_protection = "DISABLED"

    managed_scaling {
      status                    = "ENABLED"
      target_capacity           = 100
      minimum_scaling_step_size = 1
      maximum_scaling_step_size = 1
    }
  }

  tags = {
    Name        = "${local.name_prefix}-ec2-capacity-provider"
    Project     = var.project_name
    Environment = var.environment
  }
}

resource "aws_ecs_cluster_capacity_providers" "main" {
  cluster_name = aws_ecs_cluster.main.name

  capacity_providers = [
    aws_ecs_capacity_provider.ec2.name
  ]

  default_capacity_provider_strategy {
    capacity_provider = aws_ecs_capacity_provider.ec2.name
    weight            = 1
    base              = 1
  }
}

resource "aws_ecs_service" "api" {
  name            = "${local.name_prefix}-api-service"
  cluster         = aws_ecs_cluster.main.id
  task_definition = aws_ecs_task_definition.api.arn
  desired_count   = var.api_desired_count

  deployment_maximum_percent         = 200
  deployment_minimum_healthy_percent = 50
  health_check_grace_period_seconds  = 60

  capacity_provider_strategy {
    capacity_provider = aws_ecs_capacity_provider.ec2.name
    weight            = 1
    base              = 1
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.api.arn
    container_name   = "api"
    container_port   = var.api_container_port
  }

  tags = {
    Name        = "${local.name_prefix}-api-service"
    Project     = var.project_name
    Environment = var.environment
  }

  depends_on = [
    aws_ecs_cluster_capacity_providers.main,
    aws_lb_listener.api_http
  ]
}

resource "aws_iam_role" "load_generator" {
  name = "${local.name_prefix}-loadgen-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Service = "ec2.amazonaws.com"
        }
        Action = "sts:AssumeRole"
      }
    ]
  })

  tags = {
    Name        = "${local.name_prefix}-loadgen-role"
    Project     = var.project_name
    Environment = var.environment
  }
}

resource "aws_iam_role_policy_attachment" "load_generator_ssm" {
  role       = aws_iam_role.load_generator.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonSSMManagedInstanceCore"
}

resource "aws_iam_instance_profile" "load_generator" {
  name = "${local.name_prefix}-loadgen-instance-profile"
  role = aws_iam_role.load_generator.name
}

resource "aws_instance" "load_generator" {
  ami                         = data.aws_ami.amazon_linux_2023.id
  instance_type               = var.load_generator_instance_type
  subnet_id                   = aws_subnet.public[0].id
  vpc_security_group_ids      = [aws_security_group.load_generator.id]
  iam_instance_profile        = aws_iam_instance_profile.load_generator.name
  associate_public_ip_address = true
  user_data_replace_on_change = true

  user_data = <<-EOF
    #!/bin/bash
    set -eux
    dnf update -y
    dnf install -y amazon-ssm-agent docker
    systemctl enable --now amazon-ssm-agent
    systemctl enable --now docker
    docker pull alpine/bombardier:latest
  EOF

  tags = {
    Name        = "${local.name_prefix}-loadgen"
    Project     = var.project_name
    Environment = var.environment
  }

  depends_on = [
    aws_iam_role_policy_attachment.load_generator_ssm
  ]
}
