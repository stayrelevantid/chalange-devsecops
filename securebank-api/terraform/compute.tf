# Compute — dummy EC2 to fix CKV2_AWS_5 (SG not attached to resource)

resource "aws_subnet" "app" {
  vpc_id                  = aws_vpc.main.id
  cidr_block              = "10.0.1.0/24"
  availability_zone       = "${var.aws_region}a"
  map_public_ip_on_launch = false

  tags = {
    Name = "securebank-app-subnet"
  }
}

resource "aws_iam_instance_profile" "dummy" {
  name = "securebank-dummy-profile"
  role = aws_iam_role.ec2_instance.name
}

resource "aws_instance" "dummy" {
  ami                    = "ami-0c55b159cbfafe1f0"
  instance_type          = var.instance_type
  subnet_id              = aws_subnet.app.id
  vpc_security_group_ids = [aws_security_group.api.id]
  iam_instance_profile   = aws_iam_instance_profile.dummy.name
  monitoring             = true
  ebs_optimized          = true

  metadata_options {
    http_tokens                 = "required"
    http_put_response_hop_limit = 1
  }

  root_block_device {
    encrypted   = true
    kms_key_id  = aws_kms_key.securebank.arn
    volume_size = 8
    volume_type = "gp3"
  }

  tags = {
    Name = "securebank-dummy"
  }
}