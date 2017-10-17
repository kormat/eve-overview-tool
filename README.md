# eve-overview-tool

EOT is a small utility to make working with the Eve overview settings yaml
exports easier. When run against an overview yaml, it outputs a new version
with all of the Ids (group IDs, filter states, etc) annotated with their
descriptions as yaml comments.

E.g.:
```
presets:
- - pvp
  - - - alwaysShownStates
      - []
    - - filteredStates
      - - 11
        - 12
        - 14
        - 15
    - - groups
      - - 10
        - 15
        - 25
        - 26
```
becomes:
```
presets:
- - pvp
  - - - alwaysShownStates
      - []
    - - filteredStates
      - - 11 # Pilot is in your fleet
        - 12 # Pilot is in your corporation
        - 14 # Pilot is in your alliance
        - 15 # Pilot has Excellent Standing.
    - - groups
      - - 10 # Stargate
        - 15 # Station
        - 25 # Frigate
        - 26 # Cruiser
```
See a full example [here](https://gist.github.com/kormat/098d3890015f4a5a81d0cd39ea5270d7)

## Installation

EOT can be installed via:
`go get github.com/kormat/eve-overview-tool`

EOT is developed using Go 1.6, so any later Go version should work fine.

## Usage

```
eve-overview-tool -f orig.yaml > annotated.yaml
```

## Development

To rebuild `bindata.go`:
1. Install `go-bindata`:
   ```
   go get -u github.com/jteeuwen/go-bindata/...
   ```
1. Get latest Eve SDE [inventory groups file](https://www.fuzzwork.co.uk/dump/latest/invGroups.csv.bz2):
   ```
   wget https://www.fuzzwork.co.uk/dump/latest/invGroups.csv.bz2 -O data/invGroups.csv.bz2
   ```
1. Re-build `bindata.go`:
   ```
   go-bindata data/
   ```
