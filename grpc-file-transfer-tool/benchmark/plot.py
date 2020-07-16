# -*- coding: utf-8 -*-
import glob
import re

import matplotlib.pyplot as plt
import numpy as np

def average_secs(outputs):
    _sum = 0.0
    for output in outputs:
        fd = open(output)
        line = fd.readline()
        fd.close()
        ret = re.match(r'used (.*?) secs', line)
        _sum += float(ret.group(1).strip())
    return _sum / float(len(outputs))


def plot(secs_array):
    # 创建一个点数为 8 x 6 的窗口, 并设置分辨率为 80 像素/每英寸
    plt.figure(figsize=(8, 6), dpi=80)
    # 再创建一个规格为 1 x 1 的子图
    plt.subplot(1, 1, 1)
    # 柱子总数
    N = 12
    # 包含每个柱子下标的序列
    indexes = np.arange(N)
    # 包含每个柱子对应值的序列
    values = np.asarray(secs_array, dtype=np.float32)
    # 柱子的宽度
    width = 0.50
    # 绘制柱状图, 每根柱子的颜色为紫罗兰色
    plt.bar(indexes, values, width, label="transfer time", color="#87CEFA")
    # 添加数据标签
    for a, b in zip(indexes, values):
        plt.text(a, b + 0.05, '%.2f' % b, ha='center', va='bottom', fontsize=10)
    # 设置横轴标签
    plt.xlabel('chunk size')
    # 设置纵轴标签
    plt.ylabel('transfer time (s)')
    # 添加标题
    plt.title('transfer time used to send file when chunk size changed')
    # 添加纵横轴的刻度
    plt.xticks(indexes, ('1K', '2K', '4K', '8K', '16K', '32K', '64K', '128K', '256K', '512K', '1M', '2M'))
    plt.yticks(np.arange(0, 11, 1))
    # 添加图例
    plt.legend(loc="upper right")
    plt.show()


if __name__ == "__main__":
    secs_array = []
    for i in range(10, 22):
        chunk = 1<<i
        outputs = glob.glob("output_*_{}.log".format(chunk))
        secs = average_secs(outputs)
        secs_array.append(secs)
    plot(secs_array)
