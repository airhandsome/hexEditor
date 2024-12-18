# HexEditor

一个基于 Wails + Vue3 开发的十六进制编辑器。

## 功能特性

### 已实现功能
1. 基础编辑
   - [x] 文件打开/保存
   - [x] 十六进制显示和编辑
   - [x] ASCII 显示和编辑
   - [x] 撤销/重做支持
   - [x] 复制/粘贴

2. 导航功能
   - [x] 地址跳转
   - [x] 文本/十六进制查找
   - [x] 虚拟滚动
   - [x] 选择范围

3. 视图功能
   - [x] 文件面板
   - [x] 数据检查器
   - [x] 状态栏
   - [x] 工具栏

4. 数据分析
   - [x] 整数类型（8/16/32/64位）
   - [x] 浮点类型（单精度/双精度）
   - [x] 时间类型（DOS/Unix/File Time）
   - [x] 字符串类型（ASCII/UTF-16）

### 计划实现功能
1. 高级编辑
   - [ ] 块选择和编辑
   - [ ] 插入/删除字节
   - [ ] 填充数据
   - [ ] 二进制操作

2. 分析功能
   - [ ] 数据统计
   - [ ] 熵分析
   - [ ] 文件结构解析
   - [ ] 模式匹配

3. 扩展功能
   - [ ] 书签管理
   - [ ] 比较功能
   - [ ] 导出功能
   - [ ] 插件系统

## 安装和使用

### 开发环境要求
- Go 1.21+
- Node.js 18+
- Wails CLI
- Vue 3

### 安装步骤

1. 克隆仓库

```bash
git clone https://github.com/yourusername/hexeditor.git
cd hexeditor
```

2. 安装前端依赖
```bash
cd frontend
npm install
```

3. 构建并运行
```bash
cd ..
wails dev
```

### 使用方法

1. 打开文件
   - 点击工具栏的"打开"按钮
   - 使用快捷键 `Ctrl+O`
   - 拖拽文件到编辑器窗口

2. 编辑数据
   - 直接在十六进制区域或 ASCII 区域输入
   - 使用 `Tab` 键在十六进制和 ASCII 区域切换
   - 使用方向键导航

3. 查找数据
   - 使用 `Ctrl+F` 打开查找对话框
   - 支持十六进制和文本查找
   - 支持向前/向后查找

4. 数据分析
   - 使用数据检查器查看不同数据类型
   - 支持大端/小端切换
   - 支持多种时间格式

## 快捷键

| 功能 | 快捷键 |
|------|--------|
| 打开文件 | Ctrl+O |
| 保存文件 | Ctrl+S |
| 查找 | Ctrl+F |
| 跳转 | Ctrl+G |
| 撤销 | Ctrl+Z |
| 重做 | Ctrl+Y |
| 复制 | Ctrl+C |
| 粘贴 | Ctrl+V |

## 项目结构

```
hexeditor/
├── frontend/               # 前端 Vue 项目
│   ├── src/
│   │   ├── components/    # Vue 组件
│   │   ├── assets/        # 静态资源
│   │   └── main.js        # 入口文件
│   └── package.json
├── main.go                # 后端入口
└── wails.json            # Wails 配置
```

## 贡献指南

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开 Pull Request

## 许可证

[GNU License](LICENSE)

## 致谢

- [Wails](https://wails.io/)
- [Vue.js](https://vuejs.org/)
- [010 Editor](https://www.sweetscape.com/010editor/)
