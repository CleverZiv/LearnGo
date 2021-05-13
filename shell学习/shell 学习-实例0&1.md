# shell 学习-实例0&1

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

### 1.1 需求：

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

### 1.2 实现

#### 1.2.1 设置时区并同步时间

```shell
ln -s /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
if ! crontab -l |grep ntpdate &>/dev/null ; then
    (echo "* 1 * * * ntpdate time.windows.com >/dev/null 2>&1";crontab -l) |crontab 
fi
```

##### 1.2.1.1 设置时区：`ln -s /usr/share/zoneinfo/Asia/Shanghai /etc/localtime`

**`ln -s [源文件或目录][目标文件或目录]`**：在”目标目录下“放置源文件。-s 是代号（symbolic）的意思。 

当我们需要在不同的目录，用到相同的文件时，我们不需要在每一个需要的目录下都放一个必须相同的文件，我们只要在某个固定的目录，放上该文件，然后在其它的 目录下用ln命令链接（link）它就可以，不必重复的占用磁盘空间。

这里有两点要注意：

- 第一，ln命令会保持每一处链接文件的同步性，也就是说，不论你改动了哪一处，其它的文件都会发生相同的变化；
- 第二，ln的链接又软链接 和硬链接两种，软链接就是`ln -s [][]` **,它只会在你选定的位置上生成一个文件的镜像，不会占用磁盘空间，硬链接** `ln [][]`,没有参数-s, 它会在你选定的位置上生成一个和源文件大小相同的文件，无论是软链接还是硬链接，文件都保持同步变化。 
  如果你用ls察看一个目录时，发现有的文件后面有一个@的符号，那就是一个用ln命令生成的文件，用ls -l命令去察看，就可以看到显示的link的路径了。 

##### 1.2.1.2 同步时间：` (echo "* 1 * * * ntpdate time.windows.com >/dev/null 2>&1";crontab -l) |crontab `

- **`ntpdate[服务器地址]`**：`time.windows.com`是 windows 时间服务器地址，这个命令就可以将当前服务器的系统时间与windows 时间服务器的时间同步。

  > 参考：https://www.cnblogs.com/liushui-sky/p/9203657.html

- **`crontab -l`**：列出目前的任务表。

  > 参考：https://www.runoob.com/linux/linux-comm-crontab.html

- **` >/dev/null 2>&1`**：shell 中，0表示标准输入，1表示标准输出，2表示标准错误输出。
  - `>` 默认为标准输出重定向，与`1>` 相同，`2>&1` 的意思是把标准错误输出重定向到标准输出。
  - `&>file` 的意思是把标准输出和标准错误输出都重定向到文件 file 中。
  - `>/dev/null` 的意思是将标准输出重定向到文件 `/dev/null`中，这个文件比较特殊，所有传给它的东西它都丢弃掉。

- **`;`**：表示紧接着上一条命令，再执行分号后面的命令。

- `()`：表示括号中的命令作为一组命令一起按顺序执行。

- `command 1|command 2  `：管道，把第一个命令的执行结果作为第2个命令的输入传给第2个命令。

  > 参考：https://www.cnblogs.com/aaronLinux/p/8340281.html

- `crontab`：可以接收管道左边的标准输入，然后写入”任务表“中，相当于执行了 `crontab -e`后自动写入了左边的标准输入

接下来再看一下这个判断的命令:

```shell
if ! crontab -l |grep ntpdate &>/dev/null ; then
 ...
fi
```

- `if condition then ... fi`：流程控制语句。

  > 参考：https://www.runoob.com/linux/linux-shell-process-control.html

- `grep`：指令用于查找内容包含指定的范本样式的文件，这里的意思是查找包含"ntpdate"的语句。

  > 参考：https://www.runoob.com/linux/linux-comm-grep.html

整个语句的意思就是说，首先 `crontab -l` 列出当前任务表，并且筛选有关键字"ntpdate" 的内容，如果没有筛选到，则执行添加计划任务的命令

#### 1.2.2 禁用 selinux

```shell
# 禁用selinux
sed -i '/SELINUX/{s/permissive/disabled/}' /etc/selinux/config
```

**`sed`**：流编辑命令

> 参考：http://c.biancheng.net/view/4028.html

上述语句的意思是：将文件`/etc/selinux/config` 中的 `SELINUX`行中的`permissive`替换为`disable`

**选项和参数的区别：**

> 参考：http://c.biancheng.net/view/3160.html

