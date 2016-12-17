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

window.head = React.createElement("img", { src: "image/logo.gif" });

window.Navigation = React.createClass({
    displayName: "Navigation",

    render: function () {
        return React.createElement(
            Navbar,
            null,
            React.createElement(
                Navbar.Header,
                null,
                React.createElement(
                    Navbar.Brand,
                    null,
                    React.createElement(
                        "a",
                        { href: "#" },
                        "\u8BFE\u7A0B\u7F51\u7AD9"
                    )
                )
            ),
            React.createElement(
                Nav,
                null,
                React.Children.map(this.props.children, function (child) {
                    return React.createElement(
                        NavItem,
                        { href: child.props.link },
                        child
                    );
                })
            )
        );
    }
});

window.Login = React.createClass({
    displayName: "Login",

    render: function () {
        return React.createElement(
            "div",
            null,
            React.createElement(
                Well,
                { bsSize: "medium" },
                React.createElement(
                    "h2",
                    null,
                    "\u8BFE\u7A0B\u7F51\u7AD9\u767B\u5F55"
                ),
                React.createElement(
                    "form",
                    null,
                    React.createElement(
                        FormGroup,
                        { controlId: "formBasicText" },
                        React.createElement(
                            ControlLabel,
                            null,
                            "\u7528\u6237\u540D"
                        ),
                        React.createElement(FormControl, { name: "1", type: "text", placeholder: "\u8F93\u5165\u7528\u6237\u540D" }),
                        React.createElement(
                            ControlLabel,
                            null,
                            "\u5BC6\u7801"
                        ),
                        React.createElement(FormControl, { name: "2", type: "text", placeholder: "\u8F93\u5165\u5BC6\u7801" })
                    ),
                    React.createElement(
                        FormGroup,
                        null,
                        React.createElement(
                            Radio,
                            { name: "3", value: "1", inline: true, defaultChecked: true },
                            "\u6559\u5E08"
                        ),
                        React.createElement(
                            Radio,
                            { name: "3", value: "2", inline: true },
                            "\u5B66\u751F"
                        ),
                        React.createElement(
                            Radio,
                            { name: "3", value: "3", inline: true },
                            "\u52A9\u6559"
                        )
                    ),
                    React.createElement(
                        Button,
                        { type: "submit", bsSize: "large", block: true },
                        "\u767B\u5F55"
                    )
                )
            )
        );
    }
});

window.Contact = React.createClass({
    displayName: "Contact",

    render: function () {
        return React.createElement(
            Well,
            { bsSize: "medium" },
            React.createElement(
                "h2",
                null,
                "\u8054\u7CFB\u6211\u4EEC"
            ),
            React.createElement(
                ControlLabel,
                null,
                "TEL: 0571-888888"
            ),
            React.createElement(
                ControlLabel,
                null,
                "EMAIL: 666@zju.edu.cn"
            ),
            React.createElement("br", null),
            React.createElement("br", null),
            React.createElement(
                Button,
                { type: "submit", bsSize: "medium", block: true },
                "\u6211\u8981\u7559\u8A00"
            )
        );
    }
});

window.Homework = React.createClass({
    displayName: "Homework",

    render: function () {
        return React.createElement(
            Well,
            { bsSize: "medium" },
            React.createElement(
                "form",
                { method: "post", action: "http://www.baidu.com" },
                React.createElement(
                    FormGroup,
                    null,
                    React.createElement(
                        ControlLabel,
                        null,
                        "Textarea"
                    ),
                    React.createElement(FormControl, { name: "d", componentClass: "textarea", placeholder: "textarea" })
                ),
                React.createElement(
                    Label,
                    { bsStyle: "info" },
                    "\u9009\u62E9\u9644\u4EF6"
                ),
                React.createElement("br", null),
                React.createElement("br", null),
                React.createElement("input", { type: "file", name: "a", accept: "*" }),
                React.createElement("br", null),
                React.createElement(
                    Button,
                    { type: "submit" },
                    "\u63D0\u4EA4"
                )
            )
        );
    }
});

window.Forum = React.createClass({
    displayName: "Forum",

    render: function () {
        return React.createElement(
            Table,
            { striped: true, bordered: true, condensed: true, hover: true },
            React.createElement(
                "thead",
                null,
                React.createElement(
                    "tr",
                    null,
                    React.createElement(
                        "th",
                        null,
                        "\u53D1\u5E16\u8005"
                    ),
                    React.createElement(
                        "th",
                        null,
                        "\u5E16\u5B50\u6807\u9898"
                    ),
                    React.createElement(
                        "th",
                        null,
                        "\u56DE\u5E16\u6570"
                    ),
                    React.createElement(
                        "th",
                        null,
                        "\u53D1\u5E03\u65F6\u95F4"
                    )
                )
            ),
            React.createElement(
                "tbody",
                null,
                React.Children.map(this.props.children, function (child) {
                    return child;
                })
            )
        );
    }
});

window.Tiezi = React.createClass({
    displayName: "Tiezi",

    render: function () {
        return React.createElement(
            Panel,
            { header: this.props.title, bsStyle: "primary" },
            React.createElement(
                Table,
                { striped: true, bordered: true, condensed: true, hover: true, fill: true },
                React.createElement(
                    "tbody",
                    null,
                    React.Children.map(this.props.children, function (child) {
                        return child;
                    })
                )
            )
        );
    }
});

window.Respond = React.createClass({
    displayName: "Respond",

    render: function () {
        return React.createElement(
            "div",
            null,
            React.createElement(
                "form",
                { method: "", action: this.props.link },
                React.createElement(
                    FormGroup,
                    null,
                    React.createElement(
                        ControlLabel,
                        null,
                        "\u56DE\u590D\u5E16\u5B50"
                    ),
                    React.createElement(FormControl, { name: "c", componentClass: "textarea", placeholder: "textarea" })
                ),
                React.createElement(
                    Button,
                    { type: "submit" },
                    "\u53D1\u8868\u56DE\u590D"
                )
            )
        );
    }
});
