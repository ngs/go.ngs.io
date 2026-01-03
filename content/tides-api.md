---
title: tides-api
import_path: go.ngs.io/tides-api
repo_url: https://github.com/ngs/tides-api
description: A high-performance tidal prediction API written in Go, providing harmonic tidal analysis with support for multiple data sources.
version: ""
documentation_url: https://pkg.go.dev/go.ngs.io/tides-api
license: NOASSERTION
author: ngs
created_at: 2025-10-20T20:48:48Z
updated_at: 2025-10-27T09:39:20Z
---

# Tide API

A high-performance tidal prediction API written in Go, providing harmonic tidal analysis with support for multiple data sources.

## Features

- **Harmonic Tidal Analysis**: Calculate tide heights using standard tidal constituents (M2, S2, K1, O1, etc.)
- **Astronomical Nodal Corrections**: Accurate predictions using nodal corrections based on Schureman (1958)
- **Extrema Detection**: Automatically identify high and low tides with parabolic interpolation
- **Multiple Data Sources**:
  - Mock CSV data for development and testing
  - FES2014/2022 NetCDF support with bilinear interpolation
  - JMA hourly data calibration for Japanese ports
- **Flexible Configuration**: Datum offsets, timezone selection, and phase conventions
- **Clean Architecture**: Hexagonal architecture with clear separation of concerns
- **Production Ready**: Docker support, comprehensive tests, and monitoring
- **RESTful API**: Simple JSON API with ISO8601 timestamps

## Quick Start

### Prerequisites

- Go 1.22 or later
- Make (optional, for convenience commands)
- Docker (optional, for containerized deployment)

### Installation

```bash
# Clone the repository
git clone git@github.com:ngs/tides-api.git
cd tides-api

# Install dependencies
go mod download

# Copy environment configuration
cp .env.example .env

# Run the server
make run
```

The API will be available at `http://localhost:8080`.

### Using Docker

```bash
# Build Docker image
make docker-build

# Run container
make docker-run
```

## API Endpoints

### 1. Get Tide Predictions

**Endpoint**: `GET /v1/tides/predictions`

**Query Parameters**:

| Parameter | Type | Required | Description | Example |
|-----------|------|----------|-------------|---------|
| `station_id` | string | * | Station identifier | `tokyo` |
| `lat` | float | * | Latitude (-90 to 90) | `35.6762` |
| `lon` | float | * | Longitude (-180 to 180) | `139.6503` |
| `start` | string | Yes | Start time (RFC3339) | `2025-10-21T00:00:00Z` |
| `end` | string | Yes | End time (RFC3339) | `2025-10-21T12:00:00Z` |
| `interval` | string | No | Time interval (default: 10m) | `10m`, `1h` |
| `datum` | string | No | Vertical datum (default: MSL) | `MSL`, `LAT` |
| `source` | string | No | Data source (auto-detect) | `csv`, `fes` |
| `datum_offset_m` | float | No | Constant vertical offset [m] applied to all predicted heights | `0.768` |
| `timezone` | string | No | Output timezone for timestamps | `utc`, `jst` |
| `phase_convention` | string | No | Phase convention (`fes_greenwich` default, or `vu`) | `fes_greenwich`, `vu` |

\* Either `station_id` OR `lat`+`lon` must be provided (mutually exclusive)

**Example Request**:

```bash
curl 'http://localhost:8080/v1/tides/predictions?station_id=tokyo&start=2025-10-21T00:00:00Z&end=2025-10-21T12:00:00Z&interval=10m'
```

**Example Response**:

```json
{
  "source": "csv",
  "datum": "MSL",
  "timezone": "+00:00",
  "constituents": ["M2", "S2", "K1", "O1", "N2", "K2", "P1", "Q1"],
  "predictions": [
    {"time": "2025-10-21T00:00:00Z", "height_m": 0.823},
    {"time": "2025-10-21T00:10:00Z", "height_m": 0.791},
    ...
  ],
  "extrema": {
    "highs": [
      {"time": "2025-10-21T03:18:00Z", "height_m": 1.342}
    ],
    "lows": [
      {"time": "2025-10-21T09:42:00Z", "height_m": -0.187}
    ]
  },
  "meta": {
    "model": "harmonic_v0",
    "attribution": "Mock CSV (for dev). Replace with FES later."
  }
}
```

### 2. Get Constituents

**Endpoint**: `GET /v1/constituents`

Returns information about all available tidal constituents.

**Example Request**:

```bash
curl http://localhost:8080/v1/constituents
```

**Example Response**:

```json
{
  "constituents": [
    {
      "name": "M2",
      "speed_deg_per_hr": 28.9841042,
      "description": "Principal lunar semidiurnal"
    },
    ...
  ],
  "count": 18
}
```

### 3. Health Check

**Endpoint**: `GET /healthz`

Returns server health status.

**Example Request**:

```bash
curl http://localhost:8080/healthz
```

**Example Response**:

```json
{
  "status": "ok",
  "time": "2025-10-21T12:00:00Z"
}
```

## Data Sources

### CSV Mock Data (Development)

For testing without FES data, use station-based queries with mock CSV files:

1. Create a CSV file in `data/mock_{station_id}_constituents.csv`
2. Format:

```csv
constituent,amplitude_m,phase_deg
M2,0.62,145.0
S2,0.21,170.0
K1,0.18,30.0
O1,0.16,85.0
```

3. Query with `station_id`:

```bash
curl 'http://localhost:8080/v1/tides/predictions?station_id=tokyo&...'
```

See [data/README_DATA.md](data/README_DATA.md) for more details.

### FES NetCDF Data (Production) ✅ **NOW IMPLEMENTED**

For production use with FES2014/2022:

**Quick Setup:**

```bash
# 1. Install NetCDF library (macOS)
brew install netcdf

# 2. Setup AVISO credentials
make fes-setup

# 3. Download FES data
make fes-download-major  # Downloads M2, S2, K1, O1, N2, K2, P1, Q1

# 4. Start server
make run

# 5. Test with lat/lon
curl 'http://localhost:8080/v1/tides/predictions?lat=35.6762&lon=139.6503&start=2025-10-21T00:00:00Z&end=2025-10-21T12:00:00Z&interval=10m'
```

**Features:**
- ✅ Full NetCDF file reading
- ✅ Bilinear interpolation for any lat/lon
- ✅ Automatic grid caching
- ✅ Support for multiple file naming conventions
- ✅ Automatic constituent detection

**Documentation:**
- [FES_SETUP.md](FES_SETUP.md) - Complete FES setup guide
- [INSTALL.md](INSTALL.md) - Installation instructions for NetCDF library

### JMA Calibration & Station Overrides

For Japanese ports we can now calibrate directly against JMA's published hourly prediction files:

1. Download the annual TXT file for a station from `https://www.data.jma.go.jp/kaiyou/data/db/tide/suisan/txt/<year>/<station>.txt`.
2. Fit harmonic constituents using the new CLI:

```bash
# Example: fit Kisarazu (KZ) using the downloaded file
go run ./cmd/jma-harmonics \
  -jma_file /path/to/KZ.txt \
  -station KZ \
  -name "Kisarazu (KZ)" \
  -lat 35.38153 -lon 139.867951 \
  -radius_km 40 \
  > data/jma_station_overrides.json
```

3. 一括更新は Go 製ユーティリティで実行できます（TXT を `tmp/jma_txt/{CODE}.txt` に置いた上で）:

```bash
go run ./cmd/jma-overrides \
  -stations tmp/jma-stations.json \
  -txt_dir tmp/jma_txt \
  -overrides_out data/jma_station_overrides.json \
  -datum_out data/jma_datum_offsets.json
```

  `cmd/jma-overrides` は必要なら `tmp/bin/jma-harmonics` を自動ビルドし、全コード分を順次フィットします。
4. 個別に調整したい場合は `cmd/jma-harmonics` を直接叩いて JSON を追記できます。`data/jma_datum_offsets.json` も同じコマンドで併せて再生成されます。

Environment variables:

| Variable | Default | Purpose |
|----------|---------|---------|
| `DATUM_OFFSETS_PATH` | `data/jma_datum_offsets.json` | Custom path for datum offsets |
| `STATION_OVERRIDES_PATH` | `data/jma_station_overrides.json` | Custom path for constituent overrides |

With the provided Kisarazu overrides the RMSE against JMA's official hourly predictions drops below 5 cm without manual tweaking.

## Development

### Project Structure

```
tides-api/
├── cmd/
│   ├── server/              # Main API server
│   ├── jma-harmonics/       # JMA harmonic analysis tool
│   ├── jma-compare/         # JMA vs API comparison tool
│   ├── jma-overrides/       # Batch JMA station processor
│   └── fes-generator/       # FES NetCDF test data generator
├── internal/
│   ├── domain/              # Core business logic
│   │   ├── tide.go          # Tidal prediction engine
│   │   ├── constituents.go  # Standard constituent definitions
│   │   ├── nodal.go         # Astronomical nodal corrections
│   │   └── nodal_coeffs.go  # External coefficient loader
│   ├── usecase/             # Application use cases
│   │   ├── predict.go       # Prediction orchestration
│   │   └── station_adjustments.go  # JMA calibration
│   ├── adapter/             # External adapters
│   │   ├── store/           # Data stores
│   │   │   ├── csv/         # CSV mock data
│   │   │   ├── fes/         # FES NetCDF loader
│   │   │   └── bathymetry/  # GEBCO bathymetry
│   │   ├── interp/          # Bilinear interpolation
│   │   └── geoid/           # EGM2008 geoid heights
│   ├── http/                # HTTP handlers and routing
│   └── jma/                 # JMA fixed-width data parser
├── data/                    # Tidal data files
│   ├── astro_coeffs.json    # Nodal correction coefficients
│   ├── jma_datum_offsets.json      # JMA datum offsets
│   ├── jma_station_overrides.json  # JMA constituent overrides
│   └── mock_*_constituents.csv     # Mock station data
├── Makefile                 # Development commands
├── Dockerfile               # Container configuration
└── README.md
```

### Make Commands

```bash
make help              # Show all available commands
make run               # Run server locally
make build             # Build binary
make test              # Run all tests with coverage
make test-unit         # Run unit tests only
make clean             # Clean build artifacts
make fmt               # Format code
make docker-build      # Build Docker image
make docker-run        # Run in Docker
make curl-health       # Test health endpoint
make curl-tokyo        # Test Tokyo predictions
```

### Running Tests

```bash
# Run all tests with coverage
make test

# Run unit tests only (fast)
make test-unit

# Generate HTML coverage report
make test-coverage
open coverage.html
```

### Testing the API

```bash
# Start the server
make run

# In another terminal, test endpoints
make curl-health          # Health check
make curl-constituents    # List constituents
make curl-tokyo          # Tokyo predictions
make curl-tokyo-extrema  # Show high/low tides
```

## Architecture

### Domain Layer (`internal/domain/`)

Core tidal physics and calculations:

- **constituents.go**: Tidal constituent definitions and angular speeds
- **tide.go**: Harmonic analysis, tide height calculation, extrema detection

### Use Case Layer (`internal/usecase/`)

Application logic:

- **predict.go**: Orchestrates tide prediction workflow

### Adapter Layer (`internal/adapter/`)

External interfaces:

- **store/csv/**: CSV file loader for mock data
- **store/fes/**: FES NetCDF loader (stub for future implementation)
- **interp/**: Bilinear interpolation for gridded data

### HTTP Layer (`internal/http/`)

API interface:

- **handler.go**: Request handlers
- **router.go**: Route configuration

## Configuration

Environment variables (see `.env.example`):

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | Server port |
| `DATA_DIR` | `./data` | CSV data directory |
| `FES_DIR` | `./data/fes` | FES NetCDF directory |
| `GEBCO_PATH` | - | Path to GEBCO bathymetry NetCDF file |
| `MSS_PATH` | - | Path to MSS (Mean Sea Surface) NetCDF file |
| `GEOID_PATH` | - | Path to EGM2008 geoid NetCDF file |
| `ASTRO_COEFFS_PATH` | `data/astro_coeffs.json` | Path to nodal correction coefficients |
| `DATUM_OFFSETS_PATH` | `data/jma_datum_offsets.json` | Path to JMA datum offsets |
| `STATION_OVERRIDES_PATH` | `data/jma_station_overrides.json` | Path to JMA station overrides |
| `TZ` | `Asia/Tokyo` | Display timezone |

## Tidal Physics

### Harmonic Analysis

Tide height is calculated using:

```
η(t) = Σ f_k · A_k · cos(ω_k · Δt + φ_k - u_k) + MSL + datum_offset
```

Where:
- `A_k`: Amplitude of constituent k (meters)
- `φ_k`: Greenwich phase (degrees)
- `ω_k`: Angular speed (degrees/hour)
- `f_k`, `u_k`: Nodal corrections computed from astronomical arguments (Schureman 1958)
- `MSL`: Mean Sea Level offset
- `datum_offset`: Optional vertical datum adjustment (e.g., for JMA calibration)

### Supported Constituents

The API supports 18 standard tidal constituents:

**Semidiurnal** (period ~12 hours):
- M2, S2, N2, K2

**Diurnal** (period ~24 hours):
- K1, O1, P1, Q1

**Shallow Water**:
- M4, M6, MK3, S4, MN4, MS4

**Long Period**:
- Mf, Mm, Ssa, Sa

See `/v1/constituents` endpoint for full details.

### Extrema Detection

High and low tides are detected using:
1. First derivative sign change detection
2. Parabolic interpolation for sub-interval accuracy

## Implemented Features

### ✅ Completed

- [x] FES NetCDF integration with bilinear interpolation
- [x] Nodal corrections for improved accuracy (Schureman 1958)
- [x] Astronomical arguments (V0+u)
- [x] JMA hourly data calibration and harmonic fitting
- [x] Datum offset support for vertical adjustments
- [x] Custom timezone support (UTC/JST)
- [x] Phase convention options (FES Greenwich / V+u)
- [x] Automatic longitude wrapping for NetCDF grids
- [x] Complex-valued constituent support (Re/Im pairs)
- [x] Bathymetry data integration (GEBCO)
- [x] Geoid height corrections (EGM2008)

### Planned Features

- [ ] Additional vertical datums (LAT, MLLW, etc.)
- [ ] Prediction caching layer
- [ ] GraphQL API
- [ ] WebSocket streaming
- [ ] Multiple station batch queries
- [ ] OpenAPI/Swagger documentation

### Extension Points

The codebase is designed for easy extension:

- **Nodal Corrections**: External coefficient files supported via `ASTRO_COEFFS_PATH` environment variable
- **New Data Sources**: Implement `ConstituentLoader` interface in `adapter/store/store.go`
- **Custom Datums**: Use `datum_offset_m` parameter or extend `PredictionParams` in `domain/tide.go`
- **Station Overrides**: Add entries to `data/jma_station_overrides.json` for custom calibrations

## Performance

- **Latency**: <50ms for 144 points (24h @ 10min intervals)
- **Memory**: ~5MB base + ~1KB per prediction point
- **Concurrency**: Stateless design supports horizontal scaling

## License

[Specify your license here]

## Attribution

### FES Tidal Model

If using FES2014/2022 data:

> Carrère L., Lyard F., Cancet M., Guillot A. (2016). FES 2014, a new tidal model—Validation results and perspectives for improvements. In Proceedings of the ESA living planet symposium (pp. 9-13).

FES data is available from [AVISO+](https://www.aviso.altimetry.fr/) and requires registration.

### References

1. Schureman, P. (1958). Manual of Harmonic Analysis and Prediction of Tides. U.S. Coast and Geodetic Survey Special Publication No. 98. U.S. Government Printing Office, Washington, D.C.

2. Foreman, M. G. G. (1977). Manual for tidal heights analysis and prediction. Institute of Ocean Sciences, Patricia Bay.

3. Pawlowicz, R., Beardsley, B., & Lentz, S. (2002). Classical tidal harmonic analysis including error estimates in MATLAB using T_TIDE. Computers & Geosciences, 28(8), 929-937.

4. Carrère L., Lyard F., Cancet M., Guillot A. (2016). FES 2014, a new tidal model—Validation results and perspectives for improvements. In Proceedings of the ESA living planet symposium (pp. 9-13).

## Support

For issues, questions, or contributions:
- Create an issue in the repository
- Contact: Atsushi Nagase a@ngs.io

## Contributing

Contributions are welcome! Please:
1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass (`make test`)
5. Submit a pull request

## Acknowledgments

- FES team at LEGOS/CNES/CLS for tidal model data
- Japan Meteorological Agency (JMA) for hourly tide predictions data
- NOAA for EGM2008 geoid model
- GEBCO for bathymetry data
- pyTMD project for nodal correction coefficient references
- Go community for excellent libraries (Gin, go-netcdf)
