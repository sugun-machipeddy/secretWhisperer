1. Install Golang 1.16X
2. Initialize kubebuilder
    mkdir -p Projects/sercretWhisperer
    kubebuilder init --plugins go/v3 --domain shipyard.dev --owner "Sugun" --repo shipyard.dev/secretWhisperer
3. Create an API
    kubebuilder create api --group shipyard --version v1beta1 --kind Whisperer
4. Define Spec and Status of kind Whisperer
5. write the controller
6. Install CRD's on test cluster (make install)
7. Build controlller image and load into kind cluster (make docker-build IMG=secretwhisperer:1)
    (kind load docker-image secretwhisperer:1 --name thundering-typhoons)
9. Deploy the operator (make deploy IMG=secretwhisperer:1)
    or (make run)
10. Create a secret and deploy whisperer
    

