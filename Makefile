include default.config

build:
	docker build \
	  -t lambda-s3-pgn \
	  --build-arg PGN_EXTRACT_DOWNLOAD_URL=${PGN_EXTRACT_DOWNLOAD_URL} \
	  --no-cache \
	  src

build-sandbox:
	docker build \
	  -t lambda-s3-pgn-sandbox \
	  --no-cache \
	  sandbox

run-sandbox:
	docker run \
	  -e AWS_S3_BUCKET=${AWS_S3_BUCKET} \
	  -p 9000:8080 \
	  lambda-s3-pgn-sandbox

push:
	aws ecr get-login-password \
	  --region ${AWS_REGION} \
	| \
	docker login \
	  --username AWS \
	  --password-stdin \
	  ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com
	
	docker tag \
	  lambda-s3-pgn:latest \
	  ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/lambda-s3-pgn:latest
	
	docker push \
	  ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/lambda-s3-pgn:latest
