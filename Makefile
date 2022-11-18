
buildlinux:
	docker build -t crudapi --platform=linux/amd64 .
push: 
	docker tag crudapi ducnpjenkins/crudapi
	docker push ducnpjenkins/crudapi 
