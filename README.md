# ecogo
A quick go program to poke the ToU / self-powered settings on an ecoflow delta pro 3

# Building

    go build
    
# Configuring

Either set your environment variables for ACCESS_KEY and SECRET_KEY, or create a
file called ~/.ecoflow that looks like:

```
{
  "accessKey": "YOUR_ACCESS_KEY_HERE",
  "secretKey": "YOUR_SECRET_KEY_HERE",
  "serialNumber": "ECOFLOW_SERIAL_NUMBER_HERE" (optional)
}
```

# Finding your serial number

If you don't have it, run

    ./ecogo list

To get a list of Ecoflow devices associated with your account.

# Changing ToU and self-powered mode

The values the program sets are currently hard-coded to set a 70% reserve
during self-powered mode and a 45% reserve during ToU mode. You'll have to
edit ecoflow.go to change those if you want.

    ./ecogo tou      # Switch to ToU mode with a 45% backup reserve
    ./ecogo selfpow  # Self-powered mode with a 70% backup reserve

