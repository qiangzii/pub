
import os
path = "D:/Hexo_Workspace/AtuoBlog/soutce/_posts" #文件夹目录
files= os.listdir(path) #得到文件夹下的所有文件名称
s = []
for file in files: #遍历文件夹
    if not os.path.isdir(file): #判断是否是文件夹，不是文件夹才打开
        fread = open(path+"/"+file); #打开文件
        lines = f.readlines()
        fread.close()


        
