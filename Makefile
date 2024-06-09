# Makefile

generate:
	controller-gen object:headerFile="hack/boilerplate.go.txt" paths="./..."

