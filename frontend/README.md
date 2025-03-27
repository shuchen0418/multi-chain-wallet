# 多链钱包前端界面

这是多链钱包的前端界面实现，用于与后端API进行联调。

## 主要功能

- 钱包创建和导入（助记词/私钥）
- 多链支持（Ethereum、BSC、Polygon）
- 余额查询（原生代币和ERC20代币）
- 发送交易
- 交易历史和状态查询

## 技术栈

- React 18
- TypeScript
- Chakra UI (组件库)
- Axios (HTTP请求)
- Ethers.js (区块链交互)
- React Query (数据获取和缓存)

## 运行方式

1. 安装依赖
```bash
npm install
```

2. 启动开发服务器
```bash
npm run dev
```

3. 构建生产版本
```bash
npm run build
```

## 目录结构

```
frontend/
├── public/              # 静态资源
├── src/                 # 源代码
│   ├── api/             # API接口
│   ├── components/      # UI组件
│   ├── context/         # React上下文
│   ├── hooks/           # 自定义Hooks
│   ├── pages/           # 页面组件
│   ├── types/           # TypeScript类型定义
│   ├── utils/           # 工具函数
│   ├── App.tsx          # 应用入口
│   └── index.tsx        # 渲染入口
├── .env                 # 环境变量
├── package.json         # 项目依赖
└── tsconfig.json        # TypeScript配置
```

## 后端API地址配置

在`.env`文件中配置API地址：

```
REACT_APP_API_URL=http://localhost:8080/api/v1
``` 