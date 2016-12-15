window.Navbar = ReactBootstrap.Navbar; 
window.Nav = ReactBootstrap.Nav;
window.NavItem = ReactBootstrap.MenuItem;
   
window.head = (
    <img src = "image/logo.gif"/>
); //网站的抬头，这是一个浙江大学的logo


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