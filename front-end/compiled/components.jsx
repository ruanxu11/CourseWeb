var Navbar = ReactBootstrap.Navbar; 
var Nav = ReactBootstrap.Nav;
var NavItem = ReactBootstrap.MenuItem;
var Well = ReactBootstrap.Well; 
var Panel = ReactBootstrap.Panel;
var FormGroup = ReactBootstrap.FormGroup;
var ControlLabel = ReactBootstrap.ControlLabel;
var FormControl = ReactBootstrap.FormControl;
var Radio = ReactBootstrap.Radio;
var Button = ReactBootstrap.Button;
var FieldGroup = ReactBootstrap.FieldGroup;
var Label = ReactBootstrap.Label;
var Table = ReactBootstrap.Table;
var PageHeader = ReactBootstrap.PageHeader;
   
window.head = (
    <img src = "image/logo.gif"/>
); 


window.Navigation = React.createClass({
  render : function(){
    return (
        <Navbar inverse collapseOnSelect>
        <Nav>
        {
          React.Children.map(this.props.children, function (child) {
            return <NavItem href={child.props.link}>{child}</NavItem>})
        }
        </Nav>
        </Navbar>
    );
  }
});

window.Login = React.createClass({
    render: function(){
        return (
        <div>
            <Well bsSize="medium">
            <h2>课程网站登录</h2>
                <form action="/login" method="post">
                    <FormGroup controlId="formBasicText">
                    <ControlLabel>用户名</ControlLabel>
                    <FormControl name="id" type="text" placeholder="输入用户名"/>
                    <ControlLabel>密码</ControlLabel>
                    <FormControl name="password" type="text" placeholder="输入密码"/>
                    </FormGroup>
                    <FormGroup>
                    <Radio name="collection" value="teacher" inline defaultChecked>教师</Radio>
                    <Radio name="collection" value="student" inline>学生</Radio>
                    <Radio name="collection" value="teachingAssistant" inline>助教</Radio>
                    </FormGroup>
                    <Button type="submit" bsSize="large" block>登录</Button>
                </form>
                <form action="/forget/password/id" method="get"><br />
                    <Button type="submit" bsSize="medium" block>忘记密码</Button>
                </form>
            </Well>
        </div>
        );
    }
});

window.Contact = React.createClass({
    render: function(){
        return (
             <Well bsSize="medium">
                <h2>联系我们</h2>
                <ControlLabel>TEL: 0571-888888</ControlLabel>
                <ControlLabel>EMAIL: 666@zju.edu.cn</ControlLabel><br /><br />
                <Button type="submit" bsSize="medium" block>我要留言</Button>
            </Well>
        );
    }
})

window.Homework = React.createClass({
    render: function(){
        return (
            <Well bsSize="medium">
            <form method="post" action="http://www.baidu.com">
                <FormGroup>
                    <ControlLabel>Textarea</ControlLabel>
                    <FormControl name="d" componentClass="textarea" placeholder="textarea" />
                </FormGroup>
                <Label bsStyle="info">选择附件</Label><br /><br />
                <input type="file" name="a" accept="*"/><br />
                <Button type="submit">提交</Button>
            </form>
            </Well>
        );
    }
})

window.Forum = React.createClass({
    render: function(){
        return (
            <Table striped bordered condensed hover>
            <thead>
                <tr><th>发帖者</th><th>帖子标题</th><th>回帖数</th><th>发布时间</th></tr>
            </thead>
            <tbody>{React.Children.map(this.props.children, function (child) {return (child);})}</tbody>
             </Table>);        
    }
});

window.Tiezi = React.createClass({
    render: function(){
        return (
            <Panel header={this.props.title} bsStyle="primary">
                <Table striped bordered condensed hover fill>
                    <tbody>{React.Children.map(this.props.children, function (child) {return (child);})}</tbody>
                </Table>
            </Panel>
        )
    }
})

window.Respond = React.createClass({
    render: function(){
        return (
            <Well>
            <div>
            <form method="" action={this.props.link}>
                <FormGroup>
                <ControlLabel>回复帖子</ControlLabel>
                <FormControl name="c" componentClass="textarea" placeholder="textarea" />
                </FormGroup>
                <Button type="submit">发表回复</Button>
            </form>
            </div>
            </Well>
        )
    }
});

window.Introduction = React.createClass({
    render: function(){
        return (
        <div>
        <PageHeader>软件工程</PageHeader>
            <h3>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;浙江大学软件学院前身是浙江大学软件与网络学院，于2001年2月27日在杭州与宁波两地同时挂牌成立，2001年12月成为国家教育部和国家发展计划委员会批准的首批35所国家示范性软件学院之一，同时更名为浙江大学软件学院。<br/>
            &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;浙江大学国家示范性软件学院分别在杭州和宁波办学。杭州办学地点在浙江大学玉泉校区，以培养本科生为主。宁波办学地点在宁波国家高新区，以培养研究生为主。<br/><br/></h3>
        </div>
        )
    }
});

window.PanelGroup = ReactBootstrap.PanelGroup;
            window.Panel = ReactBootstrap.Panel;
            class Friend extends React.Component {
                render() {
                    return (
                        <PanelGroup defaultActiveKey="1" accordion>
                            <Panel header="友情链接" eventKey="1">
                                {this.props.links}
                            </Panel>
                        </PanelGroup>
                    );
                }
            }

            window.friendLinks = (
                <table>
                    <tr>
                        <td>
                            <a href="http://zupo.zju.edu.cn" target="_blank"><img src="http://jwbinfosys.zju.edu.cn/images/zupologo.gif" border="0" /></a>
                        </td>
                        <td width="10px"></td>
                        <td>
                            <a href="http://jwb.zju.edu.cn" target="_blank"><img src="http://jwbinfosys.zju.edu.cn/images/jwblogo.gif" border="0" /></a>
                        </td>
                        <td width="10px"></td>
                        <td>
                            <a href="http://www.cc98.org/" target="_blank"><img src="http://jwbinfosys.zju.edu.cn/images/cc98logo.gif" border="0" /></a>
                        </td>
                    </tr>
                </table>
            );

          


 