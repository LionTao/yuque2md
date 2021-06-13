# yuque2md

导出语雀成markdown文件，根据语雀的文章设置自动生成frontmatter

## Quick Start

1. 编写好`config.toml`,参照repo中的[config.toml](https://github.com/LionTao/yuque2md/blob/main/config.toml)
2. `go get github.com/liontao/yuque2md`
3. `YUQUE_TOKEN=<your token> yuque2md`

## TODO

- [ ] tags(语雀现在页面上有tags但是公开的api中无法获得，现在版本中tag默认为所在知识库的名字)