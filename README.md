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
        - 15
        - 16
    - - groups
      - - 6
        - 10
        - 15
        - 25
        - 26
        - 27
```
becomes:
```
presets:
- - pvp
  - - - alwaysShownStates
      - []
    - - filteredStates
      - - 11 # Pilot is in your fleet
        - 15 # Pilot has Excellent Standing.
        - 16 # Pilot has Good Standing.
    - - groups
      - - 6 # Celestial (2) -- Sun
        - 10 # Celestial (2) -- Stargate
        - 15 # Station (3) -- Station
        - 25 # Ship (6) -- Frigate
        - 26 # Ship (6) -- Cruiser
        - 27 # Ship (6) -- Battleship
```
See a full example [here](https://gist.github.com/kormat/098d3890015f4a5a81d0cd39ea5270d7)

## Installation

EOT can be installed via:
`go get github.com/kormat/eve-overview-tool`

EOT is developed using Go 1.10, so any later Go version should work fine.

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
1. Get latest Eve SDE [inventory category file](https://www.fuzzwork.co.uk/dump/latest/invCategories.csv.bz2):
   ```
   wget https://www.fuzzwork.co.uk/dump/latest/invCategories.csv.bz2 -O data/invCategories.csv.bz2
   ```
1. Get latest Eve SDE [inventory groups file](https://www.fuzzwork.co.uk/dump/latest/invGroups.csv.bz2):
   ```
   wget https://www.fuzzwork.co.uk/dump/latest/invGroups.csv.bz2 -O data/invGroups.csv.bz2
   ```
1. Re-build `bindata.go`:
   ```
   go-bindata data/
   ```

To update `groups/`:
1. Save the "All" overview default to a char's preset with the name "All", and export that.
1. Run `eve-overview-tool -f exported.yaml -update-groups`
