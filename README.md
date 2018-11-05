# Obot

Obot is a small program for batch processing multiple operations on bitmex api.
Information about api users is required. Api keys and secrets must be stored 
in csv format.

# Usage
Use `obot -help` to print available options.

# Providing user information
Obot can be used by:

Passing key,secret with `echo`
`echo "mi19WitPVubhsGHC2z43M8sY,mbcfenp1kRNIGShv785JACi3E1y8HV4sjoNA51-9KrFjCrfW" | obot`
or
`obot < some_data_file.csv`
or by parameter
`obot -file some_data_file.csv`
