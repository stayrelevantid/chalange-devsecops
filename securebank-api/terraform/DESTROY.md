# Destroy Plan — Emergency Teardown

## ⚠️ Kapan Dipakai?

Hanya jika `terraform apply` tidak sengaja dijalankan dan resources terbuat di AWS.

## Prasyarat

- AWS CLI terinstall + credentials configured
- Terraform terinstall
- AWS provider sudah di-download (`terraform init`)

## Langkah Destroy

### 1. Init (download provider)

```bash
terraform -chdir=securebank-api/terraform init
```

### 2. Cek apa yang akan di-destroy

```bash
terraform -chdir=securebank-api/terraform plan -destroy
```

### 3. Destroy semua resources

```bash
terraform -chdir=securebank-api/terraform destroy -auto-approve
```

### 4. Verifikasi tidak ada resources tersisa

```bash
aws s3 ls | grep securebank
aws ec2 describe-instances --query 'Reservations[].Instances[].Tags[?Key==`Name`].Value' --region ap-southeast-1
aws logs describe-log-groups --query 'logGroups[].logGroupName' --region ap-southeast-1
aws sns list-topics --region ap-southeast-1
```

### 5. Hapus state file lokal

```bash
rm -f securebank-api/terraform/terraform.tfstate*
rm -f securebank-api/terraform/.terraform.lock.hcl
rm -rf securebank-api/terraform/.terraform/
```

## Resources yang Akan Di-Destroy

| Resource | Type |
|----------|------|
| aws_vpc.main | VPC |
| aws_s3_bucket.logs | S3 |
| aws_s3_bucket.access_logs | S3 |
| aws_s3_bucket.replication_dest | S3 (us-east-1) |
| aws_security_group.api | Security Group |
| aws_default_security_group.default | Default SG |
| aws_flow_log.main | VPC Flow Log |
| aws_cloudwatch_log_group.flow_log | CloudWatch |
| aws_iam_role.flow_log | IAM Role |
| aws_iam_role.replication | IAM Role |
| aws_instance.dummy | EC2 |
| aws_subnet.app | Subnet |
| aws_sns_topic.s3_events | SNS |
| 7 S3 hardening resources | Various S3 |

## Estimasi Cost Jika Tidak Di-Destroy

- EC2 t3.micro: ~$7/bulan
- S3: ~$0.02/bulan
- CloudWatch Logs: ~$0.50/bulan
- SNS: ~$0.00 (free tier)
- **Total: ~$8/bulan jika ditinggalkan**

## Mencegah Accidental Apply

1. Jangan pernah jalankan `terraform init` di directory ini
2. Jangan set AWS credentials di environment
3. `.terraform/` dan `terraform.tfstate` sudah di `.gitignore`
4. `prevent_destroy = true` di S3 dan VPC resources (jika apply dijalankan, destroy akan gagal untuk resources ini — force delete diperlukan)