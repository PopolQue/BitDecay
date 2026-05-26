# BIT-DECAY

A rhythmic incremental/clicker game with a unique "glitch & corruption" visual aesthetic, built with **Go** and **Ebitengine**.

![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat-square&logo=go)
![Ebitengine](https://img.shields.io/badge/Ebitengine-v2.9+-FF4D4D?style=flat-square)

## 🕹️ Core Mechanics: "Rhythmic Mining"

Unlike traditional clickers, **BIT-DECAY** transforms resource generation into a high-speed rhythm game. 

- **Tempo:** 120 BPM.
- **Synchronization:** The "Manual Override" (Beat) widget rewards clicks timed to **32nd notes**.
- **Combo Multiplier:** Each rhythmic hit adds **+0.1x** to your multiplier. Missing a beat resets your combo to 1.0x.
- **Overclocking:** Maintain a perfect streak to achieve `[OVERCLOCK]` status and maximize your bit generation.

## 📺 Visual Aesthetic

The game features a custom **CRT-inspired rendering pipeline**:
- **Kage Shaders:** Real-time phosphor bloom, scanlines, and chromatic aberration.
- **Dynamic Glitch System:** As your system's `Corruption` and `Entropy` rise, the UI begins to physically decay and glitch.
- **Waterfall Background:** A procedurally generated digital "rain" that reacts to system health.

## 🛠️ Key Systems

- **Hardware Management:** Purchase and upgrade physical racks, cooling units, and power supplies.
- **Infrastructure Constraints:** Balance **Power Usage**, **Thermal Levels**, and **Rack Space**. Overloading your power grid increases `Entropy`, while overheating triggers `Corruption`.
- **System Reboot:** Prestige mechanic that purges all data for permanent `GHz Multiplier` boosts.
- **Packet Interception:** Random high-value data packets appear for manual harvesting.

## 🚀 Getting Started

### Prerequisites
- [Go 1.24+](https://golang.org/dl/)

### Running from Source
```bash
git clone https://github.com/popolque/firstbitengi.git
cd firstbitengi
go run ./cmd/bitdecay
```

### Build Commands
```bash
# Desktop (macOS)
go build -o bitdecay ./cmd/bitdecay

# Windows
GOOS=windows GOARCH=amd64 go build -o bitdecay.exe ./cmd/bitdecay
```

## ⌨️ Controls
- **Left Click:** Interact with UI, buy hardware, and hit the "Beat".
- **[Y] / [N]**: Confirm or Abort system reboot.
- **Scroll Wheel:** Scroll through Hardware and Upgrade lists.
