window.Navbar = ReactBootstrap.Navbar; 
window.Nav = ReactBootstrap.Nav;
window.NavItem = ReactBootstrap.MenuItem;
window.Well = ReactBootstrap.Well; 
window.Panel = ReactBootstrap.Panel;
window.FormGroup = ReactBootstrap.FormGroup;
window.ControlLabel = ReactBootstrap.ControlLabel;
window.FormControl = ReactBootstrap.FormControl;
window.Radio = ReactBootstrap.Radio;
window.Button = ReactBootstrap.Button;
window.FieldGroup = ReactBootstrap.FieldGroup;
var Label = ReactBootstrap.Label;
var Table = ReactBootstrap.Table;
   
window.head = (
    <img src = "image/logo.gif"/>
); 


window.Navigation = React.createClass({
  render : function(){
    return (
        <Navbar>
        <Navbar.Header>
        <Navbar.Brand>
          <a href="#">课程网站</a>
        </Navbar.Brand>
        </Navbar.Header>
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
                <form>
                    <FormGroup controlId="formBasicText">
                    <ControlLabel>用户名</ControlLabel>
                    <FormControl name="1" type="text" placeholder="输入用户名"/>
                    <ControlLabel>密码</ControlLabel>
                    <FormControl name="2" type="text" placeholder="输入密码"/>
                    </FormGroup>
                    <FormGroup>
                    <Radio name="3" value="1" inline defaultChecked>教师</Radio>
                    <Radio name="3" value="2" inline>学生</Radio>
                    <Radio name="3" value="3" inline>助教</Radio>
                    </FormGroup>
                    <Button type="submit" bsSize="large" block>登录</Button>
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
})

 