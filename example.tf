provider "hello" {}

resource "hello_world" "test1" {
  nom="test"
}

resource "hello_create" "test2" {
  create = <<EOF
  echo '{"test2":""test"}' >> test2.json
  cat test2.json
  EOF

  environment = {
    test2 = "test"
  }

}