# Quickstat
quickstat is a system monitoring tool developed in Go, designed to provide real-time information 
about the hardware performance of your system. It can output data in both standard and JSON formats, 
making it suitable for both human reading an automated processing.
## Features
- Real-time CPU and memory usage monitoring
- Output in JSON format for easy integration with other tools
- Customizable unit display for data (KB, MB, GB)
- Configurable sampling time for monitoring updates
## Installation
To install quickstat, ensure you have Go installed on your system. Then, follow these steps:
```bash
git clone https://sophuwu.site/quickstat
cd quickstat
go build -ldflags="-w -s" -trimpath
sudo install quickstat /usr/local/bin/quickstat
```
## Usage
Run the tool using the following command:
```bash
quickstat
```
### Options
- `-j`: Output in JSON format
- `-r`: Repeat output
- `-KMG`: Print in unit (K for KB, M for MB, G for GB)
- `-t<n>`: Set sampling time to n seconds
## Contributing
If you have any questions or suggestions, please open an issue.
Or feel free to fix any bugs or add new features. 
If you do, please submit a pull request.
## License
This project is licensed under the MIT License.
