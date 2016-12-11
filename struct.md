网站结构
/         GET 主页
/login    GET 登录界面 [all]
    /sutudent/3140102431        POST 提交账号密码 [学生]
    /teacher/0010633            POST 提交账号密码 [老师]
    /teachingAssistant/12345678 POST 提交账号密码 [助教]
    /admin/adminid              POST 提交账号密码 [管理员]
/logout   POST 注销 [all]
/user
    /sutudent/3140102431        GET 学生信息面板 [其他人|本人]
        classes                     本学期所上课程 [本人]
        Assignments                 未完成课程作业 [本人]
    /teacher/0010633            GET 老师信息面板 [其他人|本人]
        classes                     本学期所上课程 [本人]
        /settings
    /teachingAssistant/12345678 GET 助教信息面板 [其他人|本人]
        classes                     本学期课程 [本人]

        /update                     POST 修改信息（邮箱，电话，个人介绍） [本人]
        /password                   GET 修改密码界面 [本人]
            /oldpassword                 POST 通过旧密码修改密码 [本人]
            /security/questions          POST 通过安全问题修改密码 [本人]
        /security/questions         GET 修改安全问题界面 [本人]
            /oldpassword                 POST 通过旧密码修改安全问题 [本人]
            /security/questions          POST 通过旧安全问题修改安全问题 [本人]
/class
    /8475187ad4e2c50fed755d2ed04e1ed2 GET 课程界面 [本课程老师|本课程助教|本课程学生|其他人]
        /introduction      GET 课程简介页面
        /teaching/syllabus GET 课程大纲页面
        /announcement      GET 课程公告页面
            /update             GET&POST 更新课程公告 [本课程老师|本课程助教]
        /assignment          GET 课程作业页面 [学生|老师和助教]
            /add                GET&POST 添加作业 [老师]
            /history            GET 每次作业情况 [老师|助教]
            /assignmentid       GET 单独某项作业界面（选择，填空，大题） [学生|老师|助教]
                /submit             GET&POST 上传作业 [学生]
                /update             GET&POST 更改作业 [老师]
                /cheak              GET&POST 批作业 [老师|助教]
            /team          情况好像挺复杂的...
        /materials         GET 课程资料课件页面
            /upload             POST 上传资料 [老师|助教]
            /delete             POST 删除资料 [老师|助教]
        /forum     GET 讨论区界面 [本课程老师|本课程助教|本课程学生|其他人]
            /post/:postid POST 发言 [本课程人员]
            /delete POST 删除帖子 [老师]
        /settings          GET&POST 教学班设置界面

/course
    /21191730   GET 课程界面
        /class          POST 显示所有课程 [管理员]  
            /search          POST 搜索
            /add             POST 增加课程 [管理员]
            /add/excle       POST 增加课程excle导入 [管理员]
            /student/add
            /student/delete
            /teachingAssistant/add
            /teachingAssistant/delete
        /settings
/user
    /admin/adminid                  管理员界面
        /students          POST 显示所有学生 [管理员]  
            /search          POST 搜索
            /add             POST 增加学生 [管理员]
            /add/excle       POST 增加学生excle导入 [管理员]
        /teachers         POST 显示所有老师 [管理员]  
            /search          POST 搜索
            /add             POST 增加老师 [管理员]
            /add/excle       POST 增加老师excle导入 [管理员]
        /teachingAssistant   POST 显示所有助教 [管理员]  
            /search          POST 搜索
            /add             POST 增加助教 [管理员]
            /add/excle       POST 增加助教excle导入 [管理员]