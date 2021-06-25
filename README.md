1. Install Golang 1.16X
2. Initialize kubebuilder
    mkdir -p Projects/sercretWhisperer
    kubebuilder init --plugins go/v3 --domain shipyard.dev --owner "Sugun" --repo shipyard.dev/secretWhisperer
3. Create an API
    kubebuilder create api --group shipyard --version v1beta1 --kind Whisperer
4. Define Spec and Status of kind Whisperer
    

