#**mysql数据库用户建立流程：**
    1、CREATE USER 'username'@'host' IDENTIFIED BY 'password';  
    2、GRANT ALL ON *.* TO 'username'@'host';  
    3、flush privileges;
#mysql登录命令：
    mysql -u root -p 