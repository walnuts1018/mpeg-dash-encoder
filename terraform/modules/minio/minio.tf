resource "aws_s3_bucket" "mpeg-dash-encoder-source-upload" {
  bucket = format("mpeg-dash-encoder-source-upload%s", var.bucket_name_suffix)
}

resource "aws_s3_bucket" "mpeg-dash-encoder-output" {
  bucket = format("mpeg-dash-encoder-output%s", var.bucket_name_suffix)
}

