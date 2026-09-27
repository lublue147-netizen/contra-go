# 🎮 Contra (魂斗罗) - Go + WebAssembly 云端全自动编译与多平台部署

[![Deploy to Render](https://render.com/images/deploy-to-render-button.svg)](https://render.com/deploy?repo=https://github.com/lublue147-netizen/contra-go)
[![Vercel Live](https://img.shields.io/badge/Vercel-Online%20(Live)-success?logo=vercel)](https://contra-go.vercel.app)
[![GitHub Pages](https://img.shields.io/badge/GitHub%20Pages-Online%20(Live)-blue?logo=github)](https://lublue147-netizen.github.io/contra-go/)

本项目是用 **Go语言 (Golang)** 原生开发的经典红白机 **《魂斗罗》(Contra) 第一关群岛丛林战场**。
遵循 **“免本地构建 (Zero Local Build)”** 原则：本地无需安装配置任何 Go 编译器或构建工具，所有代码托管在 **GitHub**，由 **GitHub Actions** 自动进行云端交叉编译为 WebAssembly，并全自动流水线部署到 **Vercel** 和 **Render**。

---

## 🌟 游戏核心特性

- **纯正 NES 像素风格与 60FPS 流畅运行**：采用标准 Go 1.22 + WebAssembly (`syscall/js`) 直接驱动 HTML5 高性能 Canvas 2D。
- **全套经典武器系统**：
  - `普通步枪 (Normal)`：基础白黄色单发点射。
  - `经典散弹枪 [S]`：标志性红光5向扇形散射，全屏制霸！
  - `激光枪 [L]`：高穿透青色能量光束。
  - `重机枪 [M]`：高速连续压制弹幕。
  - `无敌防护罩 [B]`：周身能量球环绕，触碰敌人直接秒杀。
- **全功能主角动作系统 (Bill Rizer)**：
  - 8 方向全方位射击与奔跑（前、上、下、斜上45°、斜下45°）。
  - 卧倒匍匐躲避子弹射击。
  - 经典空中翻滚跳跃（Somersault Jump）。
  - 下跳穿越跳跃（`↓` + `跳跃` 穿越下落悬空木桥）。
  - 水域半身潜水游泳与射击。
- **丰富的敌军与第一关要塞 Boss**：
  - 红衣奔跑突击步兵、高地掩体狙击手、360°旋转地面炮台、空中正弦波飞行的红鹰武器胶囊。
  - **第一关终极Boss【格尔马防御要塞 (Gorma Fortress)】**：顶部榴弹炮、两侧旋转防空炮塔、步兵出击舱门、核心跳动红鹰弱点、阶段破坏与全屏连锁爆炸！
- **双主角自由切换 (Bill Rizer / Lance Bean)**：
  - `比尔·雷泽 (1P)`：蓝裤、红头带、棕发。
  - `兰斯·比恩 (2P)`：红裤、蓝头带、黑发。
  - 按 <kbd>C</kbd> 键或屏幕按钮随时一键切换操控角色！
- **第一关经典索桥连锁塌陷冲刺**：
  - 踏上丛林悬空木桥后，桥面逐段点燃并连锁爆炸，逼真还原红白机生死冲锋时刻！
- **爆炸屏幕震动 (Screen Shake)**：
  - 散弹射击、桥梁引爆与 Boss 破损爆炸时带有真实街机抖动打击感。
- **纯代码 8-bit Web Audio 音乐与音效合成器**：
  - 纯代码算法实时合成热血澎湃的 **《魂斗罗》第一关群岛丛林经典 BGM**（包含方波低音、主旋律音轨与白噪声鼓点），以及全套经典射击爆破音效，无需下载外部 MP3/WAV，加载零延迟！


---

## 📁 项目目录结构

```text
├── .github/
│   └── workflows/
│       └── build-and-deploy.yml   # GitHub Actions: 自动拉取Go编译为Wasm，并分发到Vercel与Render
├── cmd/
│   └── wasm/
│       ├── main.go                # Wasm 入口、主循环 requestAnimationFrame 与事件监听
│       ├── game.go                # 游戏状态机、输入同步、全局碰撞结算
│       ├── player.go              # 玩家物理、动作状态、8向瞄准、30命秘籍
│       ├── weapon.go              # 弹道动力学、穿透计算与掉落武器拾取
│       ├── enemy.go               # 步兵/狙击手/炮台/红鹰飞行胶囊 AI 与刷怪器
│       ├── boss.go                # 第一关格尔马防御要塞 Boss 核心与炮塔
│       ├── stage.go               # 关卡平台地形、落水判定、视差滚动与镜头锁定
│       ├── audio.go               # 纯代码 8-bit Web Audio 音效合成器
│       ├── render.go              # 复古像素渲染引擎与 HUD 界面
│       └── types.go               # 核心结构体与常量定义
├── web/
│   ├── index.html                 # 街机复古外框、Wasm 加载器、虚拟按键
│   ├── style.css                  # 响应式布局、CRT 扫描线着色器与霓虹 UI
│   └── wasm_exec.js               # Go 官方 WebAssembly 运行时胶水代码
├── build.sh                       # 通用云端构建脚本（GitHub/Vercel/Render 均可调用）
├── vercel.json                    # Vercel 静态托管与 wasm mime-type 规则
├── render.yaml                    # Render 官方一键部署 Blueprint 配置
├── render-build.sh                # Render 自动化云构建脚本
├── server.go                      # 备用 Go HTTP 服务器（用于 Docker 或 Render Web Service）
├── Dockerfile                     # 多阶段 Docker 容器构建文件
├── go.mod                         # Go 模块配置
└── README.md
```

---

## 🚀 部署指南（免本地构建，全云端执行）

### 第一步：将代码推送到 GitHub

1. 在 [GitHub](https://github.com/new) 上创建一个新的仓库（例如 `contra-go`）。
2. 在本地执行如下命令，将代码提交并推送到你的 GitHub 仓库（本地无需安装任何编译环境）：

```bash
git init
git add .
git commit -m "feat: initial commit contra game in Go Wasm"
git branch -M main
git remote add origin https://github.com/<你的用户名>/<你的仓库名>.git
git push -u origin main
```

推送成功后，GitHub Actions 将会自动触发并执行远程编译！

---

### 第二步：部署到 Vercel

本项目提供 **两种** 完全免本地构建的 Vercel 部署方式：

#### 方式 A：GitHub 仓库直接一键导入（最简单推荐）
1. 打开 [Vercel 仪表盘](https://vercel.com/dashboard)，点击 **"Add New..." -> "Project"**。
2. 导入刚才推送的 GitHub 仓库。
3. Vercel 会自动识别项目根目录中的 `vercel.json` 和 `build.sh`。
4. 无需更改任何配置，直接点击 **"Deploy"**。
5. Vercel 的云端构建容器会自动下载 Go 编译器，生成 `dist/contra.wasm`，并在 1 分钟内完成全球 CDN 部署！

#### 方式 B：通过 GitHub Actions 流水线自动推送
如果希望由 GitHub Actions 编译完成后直接推送发布到 Vercel：
1. 在 Vercel 控制台获取 Token：[Vercel Tokens](https://vercel.com/account/tokens)。
2. 在 GitHub 仓库设置中，进入 **Settings -> Secrets and variables -> Actions**，添加以下三个机密：
   - `VERCEL_TOKEN`：你的 Vercel 访问令牌。
   - `VERCEL_ORG_ID`：你的团队/个人 ID（可在本地运行 `npx vercel link` 或项目设置中查看）。
   - `VERCEL_PROJECT_ID`：你的 Vercel 项目 ID。
3. 每次 `git push` 到 `main` 分支时，GitHub Actions 会编译 Wasm 并直接发布到 Vercel 生产环境。

---

### 第三步：部署到 Render

本项目已配置官方规范的 `render.yaml` Blueprint 与 `render-build.sh`，同样支持两种免本地构建方式：

#### 方式 A：Render 官方蓝图一键部署（推荐）
1. 登录 [Render 控制台](https://dashboard.render.com/)。
2. 点击右上角 **"New +" -> "Blueprint"**。
3. 连接你的 GitHub 账号并选择 `contra-go` 仓库。
4. Render 会自动读取仓库中的 `render.yaml`，自动创建 **Static Site**。
5. 点击 **"Apply"**，Render 将在云端自动拉取代码、调用 `render-build.sh` 编译 WebAssembly，并部署上线！

#### 方式 B：使用 Render Deploy Hook 配合 GitHub Actions
1. 在 Render 创建好 Static Site 后，进入 **Settings -> Deploy Hook**。
2. 复制生成的 Webhook URL（形如 `https://api.render.com/deploy/srv-xxxx?key=yyyy`）。
3. 在 GitHub 仓库的 **Settings -> Secrets and variables -> Actions** 中添加：
   - `RENDER_DEPLOY_HOOK_URL`：粘贴你的 Deploy Hook Webhook URL。
4. 之后每次 GitHub Actions 编译成功，就会自动向该 Webhook 发送触发请求，Render 会自动同步发布最新版本！

---

### 🎁 附赠福利：GitHub Pages 自动上线
GitHub Actions 工作流中还贴心集成了 GitHub Pages 自动部署。只要你在 GitHub 仓库的 **Settings -> Pages** 中将部署来源设置为 **GitHub Actions** 或 `gh-pages` 分支，推送到 GitHub 后还会立即获得一个免费可访问的：
`https://<你的用户名>.github.io/<仓库名>/`

---

## 🎮 详细操作指南

| 功能操作 | 键盘按键 (PC) | 手柄按键 (Gamepad) | 触控屏幕 (手机/平板) |
| :--- | :--- | :--- | :--- |
| **向左移动** | `A` 或 `←` | 摇杆左 / 十字键左 | 虚拟十字键 ◀ |
| **向右移动** | `D` 或 `→` | 摇杆右 / 十字键右 | 虚拟十字键 ▶ |
| **向上瞄准** | `W` 或 `↑` | 摇杆上 / 十字键上 | 虚拟十字键 ▲ |
| **卧倒匍匐** | `S` 或 `↓` | 摇杆下 / 十字键下 | 虚拟十字键 ▼ |
| **斜角瞄准** | `W+D` / `W+A` / `S+D` / `S+A` | 摇杆斜向 | 虚拟十字键斜推 |
| **射击开火** | `J` / `Z` / `Space 空格` | `B键` / `X键` (按键1/2) | 虚拟 B 键 (射击) |
| **跳跃 (空翻)**| `K` / `X` | `A键` (按键0) | 虚拟 A 键 (跳跃) |
| **下落跳跃** | `S` / `↓` + 跳跃按键 | 摇杆下 + `A键` | 虚拟十字键 ▼ + A 键 |
| **开始 / 暂停** | `Enter 回车` | `Start 键` (按键9) | START 按钮 |
| **经典 30 命** | `↑↑↓↓←→←→BA` | `↑↑↓↓←→←→BA` | 点击 **★ 30命秘籍** 按钮 |
| **音效开关** | `M` 键 | - | 顶部 🔊 音效 按钮 |
| **CRT 扫描线**| - | - | 顶部 📺 CRT 滤镜 按钮 |
| **全屏模式** | - | - | 顶部 ⛶ 全屏 按钮 |

---

## 🛡️ 武器道具详解

击破在空中正弦波滑翔的 **红色猎鹰胶囊** 或特定地面感应器，即可掉落徽章：

| 图标 | 武器名称 | 特色介绍 |
| :---: | :---: | :--- |
| **`S`** | **散弹枪 (Spread Gun)** | **魂斗罗最强标志！** 每次射击向前方同时射出 5 颗扇形红光弹丸，覆盖大半屏幕，清理杂兵和 Boss 的首选神兵！ |
| **`L`** | **激光枪 (Laser Gun)** | 超高速青色能量光束，具备穿透能力，可一次贯穿多名敌人或造成多段连续暴击！ |
| **`M`** | **重机枪 (Machine Gun)** | 极速全自动连射，长按开火即可形成密不透风的黄色子弹风暴！ |
| **`B`** | **防护罩 (Barrier)** | 激活环绕主角旋转的能量力场，10 秒内无视任何敌弹攻击，触碰任何敌军即可将其瞬间粉碎！ |

---

## 💻 本地预览说明（可选）

如需要在本地电脑进行极速预览（不需要本地 Go 编译器，如果系统有 Python 或 Node.js 也可直接预览 `dist`）：

```bash
# 使用 Python 启动测试服务器预览
python3 -m http.server 8080 --directory dist

# 或者如果有 Go 环境，也可以运行自带的极速服务器：
go run server.go
```
访问 `http://localhost:8080` 即可开始游戏。
