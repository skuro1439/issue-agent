# for Developer

## Development
### make image
After changing code, build container image
```sh
make image/dev
```

### Run
```sh
GITHUB_TOKEN=$(gh auth token) \
ANTHROPIC_API_KEY="key" \
OPENAI_API_KEY="key" \
  go run cmd/runner/main.go issue \
     --base_branch main \
     --github_issue_number 10 \
     --work_repository hobby \
     --github_owner clover0 \
     --model "gpt-4o" \
     --language Japanese \
     --log_level debug \
     --aws_region us-east-1
```

## Run using release image
```sh
# build image using Dockerfile on actual release
make image/release

# run 
GITHUB_TOKEN=$(gh auth token) \
ANTHROPIC_API_KEY="key" \
OPENAI_API_KEY="key" \
  go run -ldflags "-X main.containerImageTag=dev-release" cmd/runner/main.go issue \
     --base_branch main \
     --github_issue_number 10 \
     --work_repository hobby \
     --github_owner clover0 \
     --model "gpt-4o" \
     --language Japanese \
     --log_level debug \
     --aws_region us-east-1
```
