curl -L -s https://api.github.com/repos/silverhairs/lox/releases/latest \
| grep "browser_download_url.*glox-$(uname -s | tr '[:upper:]' '[:lower:]' | sed 's/darwin/macos/')-$(uname -m | sed 's/x86_64/amd64/')" \
| cut -d '"' -f 4 \
| wget -qi - \
&& chmod +x glox-$(uname -s | tr '[:upper:]' '[:lower:]' | sed 's/darwin/macos/')-$(uname -m | sed 's/x86_64/amd64/') \
&& mv glox-$(uname -s | tr '[:upper:]' '[:lower:]' | sed 's/darwin/macos/')-$(uname -m | sed 's/x86_64/amd64/') $PWD/glox
