# 🚀 djgo — Blazing-Fast Django Scaffolder

**djgo** is a Go-powered CLI tool that eliminates the boring, repetitive setup work every new Django project needs. No more hand-editing `urls.py`, tweaking `settings.py`, or copy-pasting Docker configs. One command in your terminal — and your project architecture is ready to go.

![Go](https://img.shields.io/badge/Go-1.2x-00ADD8?logo=go&logoColor=white)
![Django](https://img.shields.io/badge/Django-ready-092E20?logo=django&logoColor=white)
![Status](https://img.shields.io/badge/status-in%20development-yellow)
![License](https://img.shields.io/badge/license-MIT-blue)

---

## ✨ Features

- 🧩 **Full automation** — new apps are automatically registered in `INSTALLED_APPS`.
- 🧭 **Smart routing** — generates `urls.py` for every app and wires them into the main router.
- 🐍 **Isolated environment** — creates a `venv`, installs dependencies, and produces a clean `pip freeze`.
- 🐳 **Infrastructure out of the box** — generates a ready-to-use `Dockerfile`, `docker-compose.yml`, and `.gitignore`.
- 🎬 **Cinematic launch** — automatically opens the generated project in your IDE (PyCharm, VS Code, Cursor supported).

---

## 📦 Installation

The project is still under active development, so for now you'll need to build it from source:

```bash
git clone https://github.com/AntonUniatitskyi/djgo.git
cd djgo
go build -o djgo.exe main.go
```

> 💡 **Tip:** Add the folder containing the binary to your global `PATH` so you can run `djgo` from anywhere on your system.

---

## 🛠 Usage

**Generate a modern (flat) project structure with apps and Docker:**

```bash
djgo init MyProject --apps users,api,blog --docker
```

**Generate a classic Django (nested) structure:**

```bash
djgo init MyProject --apps catalog --duplicate
```

---

## 🏗 What happens under the hood?

1. 🔧 A virtual environment is created and a fresh Django install is set up.
2. 🏛 The project core is generated.
3. 📂 Folders for every app are created, each with a ready-made `urls.py`.
4. 🐳 Docker configs are configured automatically.
5. 🚀 The project opens automatically in your favorite IDE.

---

## 🤝 Contributing

Contributions, issues, and feature requests are welcome! Feel free to check the [issues page](https://github.com/YOUR_USERNAME/djgo/issues) or open a PR.

## 📄 License

This project is licensed under the MIT License.
