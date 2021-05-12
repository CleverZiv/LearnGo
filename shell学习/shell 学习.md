# shell 学习

## 0. 实例0

### 需求

输入：一个 csv 文件，其中存储了两个字段的映射关系，希望批量生成更新 sql

### 实现

```shell
#! /bin/bash
out='pd.sql' 
IFS=','
while read commodity_code outer_id
do 
	cat << end >> $out
	update t_product set commodity_id = (
		select commodity_id
		from t_standard_commodity
		where commodity_code = '$commodity_code')
	where outer_id = $outer_id;
end
done
```

其中 ”end“ 是自定义的EOF（End Of File），也即在文件中硬编码了

使用：

```shell
./test.sh < pd.csv
```



## 1. 实例1-服务器系统配置初始化

### 需求：

新购买10台服务器并已安装 Linux 操作系统

1. 设置时区并同步时间
2. 禁用 selinux
3. 清空防火墙默认策略
4. 历史命令显示操作时间
5. 禁止 root 远程登录
6. 禁止定时任务发送邮件
7. 设置最大打开文件数
8. 减少 swap 使用（swap会影响性能）
9. 系统内核参数优化（TCP连接等）
10. 安装系统性能分析工具及其它

### 实现

#### 设置时区并同步时间

```shell
ln -s /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
if ! crontab -l |grep ntpdate &>/dev/null ; then
    (echo "* 1 * * * ntpdate time.windows.com >/dev/null 2>&1";crontab -l) |crontab 
fi
```

`ln -s [源文件或目录][目标文件或目录]`：在”目标目录下“放置源文件。-s 是代号（symbolic）的意思。 

当我们需要在不同的目录，用到相同的文件时，我们不需要在每一个需要的目录下都放一个必须相同的文件，我们只要在某个固定的目录，放上该文件，然后在其它的 目录下用ln命令链接（link）它就可以，不必重复的占用磁盘空间。

这里有两点要注意：

- 第一，ln命令会保持每一处链接文件的同步性，也就是说，不论你改动了哪一处，其它的文件都会发生相同的变化；
- 第二，ln的链接又软链接 和硬链接两种，软链接就是`ln -s [][]` **,它只会在你选定的位置上生成一个文件的镜像，不会占用磁盘空间，硬链接** `ln [][]`,没有参数-s, 它会在你选定的位置上生成一个和源文件大小相同的文件，无论是软链接还是硬链接，文件都保持同步变化。 
  如果你用ls察看一个目录时，发现有的文件后面有一个@的符号，那就是一个用ln命令生成的文件，用ls -l命令去察看，就可以看到显示的link的路径了。 