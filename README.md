# Free Hands Onmyoji
>解放你的双手，不要让你在游戏中也要工作!  
阴阳师的自动化脚本，旨在帮助玩家自动执行游戏中的重复任务。
## 介绍
使用Golang和OpenCV+robotgo实现的自动化脚本，旨在帮助玩家自动执行《阴阳师》中的重复任务，如探索、寻怪等。该项目主要针对第28章的探索任务，使用模板匹配技术识别游戏界面元素，并通过模拟点击操作来完成任务。

### 主要特性
- 🎮 支持多种游戏模式：探索、突破、业园火等
- ⏰ 智能退出机制：支持手动退出和定时自动退出
- 🖥️ 多显示器支持：可在主显示器或扩展显示器上运行
- 🔄 自动化管理：可选择在退出时自动关闭BlueStacks
- 📱 模拟器兼容：专为BlueStacks模拟器优化
- 🛠️ 灵活配置：丰富的命令行参数满足不同需求

## 灵感来源
本项目受[auto-click](https://github.com/WinterBokeh/auto-click)启发

## 免责声明

### 使用风险
本项目（Free-Hands-Onmyoji）是一个自动化工具，仅用于学习和研究目的。使用本工具可能违反游戏《阴阳师》的用户协议或服务条款。使用者需自行承担因使用本工具而可能导致的所有风险，包括但不限于账号封禁、游戏数据丢失等后果。

### 非官方性质
本项目与《阴阳师》游戏开发商网易游戏没有任何关联，不是官方认可的工具。这是一个由玩家创建的第三方工具。

### 无担保
本软件按"原样"提供，不提供任何形式的明示或暗示的保证，包括但不限于对适销性、特定用途适用性和非侵权性的保证。作者或版权持有人在任何情况下均不对任何索赔、损害或其他责任负责，无论是在合同诉讼、侵权行为或其他方面。

### 合理使用
我们强烈建议用户：
1. 仅将本工具用于个人学习和研究
2. 不要长时间持续使用自动化工具
3. 遵守游戏开发商的用户协议和服务条款
4. 尊重游戏开发者的知识产权和劳动成果

使用本工具即表示您已阅读、理解并同意上述免责声明的所有条款。


## 开始运行
```shell
make run
```
### 前置条件

- Go 1.24 or later
-  opencv 4.11.0_1
-  macOS M系列

### 安装

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/free-hands-onmyoji.git
   cd free-hands-onmyoji
   ```

2. Install dependencies:
   ```
   go mod tidy
   ```
3. build
   ```
   make build
   ```
### 使用方法

  - 打开BlueStacks模拟器并且进入游戏探索具体看以下图片.  
  
   ![进入游戏](./document/20250622194936.jpg)
- 启动应用程序:
  ```bash
  # 基本使用
  ./free-hands-onmyoji -task k28
  
  # 使用扩展显示器
  ./free-hands-onmyoji -task breaker -display 1
  
  # 设置60分钟后自动退出
  ./free-hands-onmyoji -task k28 -timeout 60
  
  # 定时退出时自动关闭BlueStacks
  ./free-hands-onmyoji -task breaker -timeout 30 -closeBlueStacks
  ```
  程序会自动点击章节并进入探索界面，开始自动执行任务。

### 命令行参数

- **任务类型** (`-task`):
  - `k28`: 执行第28章探索任务
  - `guren`: 执行业园火任务
  - `breaker`: 执行突破任务（默认）
  - `limitedEvents`: 执行周期性活动
    > 周期性任务需要自己手动进入活动界面并且对需要点击的图片进行截图放到`./limitedEvents`目录下命名规则自行查看文件夹下的图片名称

- **显示器选择** (`-display`):
  - `-1`: 主显示器（默认）
  - `1`: 扩展显示器

- **定时退出** (`-timeout`):
  - `0`: 不限制运行时间（默认）
  - `>0`: 运行指定分钟数后自动退出

- **自动关闭BlueStacks** (`-closeBlueStacks`):
  - `false`: 定时退出时不关闭BlueStacks（默认）
  - `true`: 定时退出时自动关闭BlueStacks
  - **注意**: 手动退出（Command+Shift+O）始终不会关闭BlueStacks

### 退出方式

1. **手动退出**: 按下 `Command+Shift+O` 组合键
   - 立即停止程序
   - 保持BlueStacks运行，方便继续游戏

2. **定时退出**: 设置 `-timeout` 参数
   - 运行指定时间后自动退出
   - 可选择是否关闭BlueStacks（`-closeBlueStacks`参数）

### 使用场景示例

```bash
# 场景1：短时间挂机，随时可能回来游戏
./free-hands-onmyoji -task k28 -timeout 30
# 30分钟后自动退出，保持BlueStacks运行

# 场景2：长时间无人值守，完成后彻底关闭
./free-hands-onmyoji -task breaker -timeout 120 -closeBlueStacks  
# 2小时后自动退出并关闭BlueStacks

# 场景3：使用外接显示器游戏
./free-hands-onmyoji -task guren -display 1
# 在扩展显示器上运行，手动退出

# 场景4：测试脚本，随时手动停止
./free-hands-onmyoji -task k28
# 无时间限制，按Command+Shift+O随时退出
```

### 注意事项
- 启动后请确保BlueStacks窗口处于活动状态。
- 脚本运行过程中请勿手动操作BlueStacks窗口，以免干扰脚本执行。
- 如果脚本无法正确识别界面元素，可能需要调整截图或更新模板图片。
- 目前仅支持macOS M系列，其他平台可能需要额外的适配工作。
- 目前仅支持BlueStacks模拟器，其他模拟器可能需要额外的适配工作。
- 使用定时退出功能时，建议设置合理的时间长度，避免过度使用自动化工具。
- 在多显示器环境下，请确保BlueStacks在正确的显示器上运行。
- 自动关闭BlueStacks功能会强制结束BlueStacks进程，请确保游戏数据已保存。


### FAQ

- **为什么需要使用BlueStacks模拟器？**
  - 因为该脚本依赖BlueStacks只针对了BlueStacks来获取窗口位置

- **为什么需要使用macOS M系列？**
  - 作者是基于macOS M系列进行开发和测试的，其他平台可能需要额外的适配工作。

- **为什么需要OpenCV？**
  - OpenCV用于图像处理和模板匹配，以识别游戏中的元素。

- **匹配不到图片怎么办？**
  - 使用截图工具重新截图 因为位置计算和模版图片在不同的窗口大小下有区别所以需要重新截图 项目中的图片是基于作者自己的显示器+MacOS的自动4分窗口截取不是适合所有人

- **构建报错怎么办？**
  - 可能是openCV路径不正确或者是GCC配置不对 需要你自行修改Makefile中的路径

- **如何设置定时运行？**
  - 使用 `-timeout` 参数设置运行时长，例如：`./free-hands-onmyoji -task k28 -timeout 60` 表示运行60分钟后自动退出

- **手动退出和定时退出有什么区别？**
  - **手动退出** (`Command+Shift+O`)：立即停止脚本，保持BlueStacks运行
  - **定时退出** (`-timeout`)：可选择是否关闭BlueStacks（`-closeBlueStacks` 参数）

- **如何在双屏环境下使用？**
  - 使用 `-display 1` 参数选择扩展显示器，默认使用主显示器
  - openCV的安装
    - 如果你使用的是macOS M系列可以直接使用brew安装
      ```shell
      brew install opencv
      ```
    - 如果你使用的是其他平台请参考[OpenCV官方文档](https://opencv.org/releases/)进行安装

## Todo

- [x] 困28章的探索脚本
- [x] 优化日志输出
- [x] 增加定时退出功能
- [x] 增加BlueStacks自动关闭功能
- [x] 支持多显示器环境
- [x] 区分手动退出和定时退出行为
- [ ] 增加更多的测试用例
- [ ] 增加御魂探索脚本
- [x] 增加突破脚本
- [x] 增加御灵探索脚本
- [x] 增加业园火脚本
## License

本项目采用MIT许可证授权，这意味着您可以自由地使用、修改和分发本软件，无论是用于个人还是商业用途，只要您保留原始版权声明。

详细条款请参阅[LICENSE](./LICENSE)文件。

```
MIT License

Copyright (c) 2025 Free-Hands-Onmyoji Contributors

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files...
```