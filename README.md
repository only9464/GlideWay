![GlideWay](https://socialify.git.ci/only9464/GlideWay/image?description=1&font=Inter&forks=1&issues=1&logo=https%3A%2F%2Fglideway.github.io%2FDevelopmentDocs%2Fimages%2FG.png&name=1&owner=1&pattern=Floating%20Cogs&pulls=1&stargazers=1&theme=Auto)
[![go](https://img.shields.io/github/go-mod/go-version/only9464/GlideWay)]()
[![size](https://img.shields.io/github/repo-size/only9464/GlideWay)]()
[![commit](https://img.shields.io/github/last-commit/only9464/GlideWay)]()
[![contributors](https://img.shields.io/github/contributors/only9464/GlideWay)]()
[![stable](https://img.shields.io/badge/stable-stable-green.svg)](https://github.com/only9464/GlideWay)

<!-- [![license](https://img.shields.io/github/license/only9464/GlideWay)]()
[![download](https://img.shields.io/github/downloads/only9464/GlideWay/total)]()
[![release](https://img.shields.io/github/v/release/only9464/GlideWay)]() -->

<!-- <p align="center">
<img src="https://img.shields.io/badge/wails-v2.9.2-red">
</p> -->

---

## 声明

**本程序仅供于学习交流，请使用者遵守《中华人民共和国网络安全法》，勿将此工具用于非授权的测试，开发者不负任何连带法律责任。**

## 介绍

- [X] 端口扫描器
  - [X] TCP快速扫描
  - [X] 线程池扫描
  - [X] 端口服务类型、产品名称、版本信息
  - [X] 操作系统、主机名、设备类型、探针名称
  - [ ] 可能存在的漏洞
  - [ ] CVE漏洞针对性扫描
  - [ ] 结果导出
- [X] 目录扫描器
  - [X] 自定义字典文件
  - [X] 线程池扫描
  - [X] 响应状态码
  - [X] 响应类型
  - [X] 响应体大小
  - [X] 结果多功能排序
  - [ ] Cookie扫描
  - [ ] 自定义User-Agent
  - [ ] 指纹识别
  - [ ] 可能存在的漏洞
  - [ ] CVE漏洞扫描
  - [ ] 结果导出
- [X] Gitdorker
  - [X] 自定义主关键词
  - [X] 自定义副关键词(可字典文件)
  - [X] 字典文件多种方式切割
  - [X] Github Web Search Url 拼接(一键跳转)
  - [X] 结果精准定位到仓库源码文件
  - [ ] 适配 Gitlab (包括自部署)
  - [ ] 结果导出
- [ ] JSFinder
  - [ ] 多线程扫描
  - [ ] 敏感词结果整合
  - [ ] 加解密算法整合导出
  - [ ] 完整路径拼接
  - [ ] 子域名收集
  - [ ] 深度扫描
  - [ ] 结果导出
- [ ] SQLScan
  - [ ] 自定义字典文件
  - [ ] 多线程扫描
  - [ ] Cookie扫描
  - [ ] 自定义User-Agent
  - [ ] 结果导出

## 关于

本项目为哈尔滨工程大学计算机科学与技术学院2024~2025学年秋季学期**数据安全课程设计**课程的课程设计，原选题为**基于Wails框架的网络扫描工具**

## 鸣谢

下列是本项目使用或者参考的优秀开源框架，感谢网上众多的开源项目及其开源项目的作者，致敬为网络安全事业做出贡献的每一位前辈！(排名不分先后)

* [Slack](https://github.com/qiwentaidi/Slack) - 安全服务集成化工具平台，希望能帮助你少开几个应用测试
* [GitDorker](https://github.com/obheda12/GitDorker) - A Python program to scrape secrets from GitHub through usage of a large repository of dorks.
* [gonmap](https://github.com/lcvvvv/gonmap) - gonmap是一个go语言的nmap端口扫描库，使用纯go实现nmap的扫描逻辑，而非调用nmap来进行扫描。
* [gobuster](https://github.com/OJ/gobuster) - Directory/File, DNS and VHost busting tool written in Godetection