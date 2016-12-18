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

window.head = React.createElement("img", { src: "image/logo.gif" });

window.Navigation = React.createClass({
    displayName: "Navigation",

    render: function () {
        return React.createElement(
            Navbar,
            { inverse: true, collapseOnSelect: true },
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
                    { action: "/login", method: "post" },
                    React.createElement(
                        FormGroup,
                        { controlId: "formBasicText" },
                        React.createElement(
                            ControlLabel,
                            null,
                            "\u7528\u6237\u540D"
                        ),
                        React.createElement(FormControl, { name: "id", type: "text", placeholder: "\u8F93\u5165\u7528\u6237\u540D" }),
                        React.createElement(
                            ControlLabel,
                            null,
                            "\u5BC6\u7801"
                        ),
                        React.createElement(FormControl, { name: "password", type: "text", placeholder: "\u8F93\u5165\u5BC6\u7801" })
                    ),
                    React.createElement(
                        FormGroup,
                        null,
                        React.createElement(
                            Radio,
                            { name: "collection", value: "teacher", inline: true, defaultChecked: true },
                            "\u6559\u5E08"
                        ),
                        React.createElement(
                            Radio,
                            { name: "collection", value: "student", inline: true },
                            "\u5B66\u751F"
                        ),
                        React.createElement(
                            Radio,
                            { name: "collection", value: "teachingAssistant", inline: true },
                            "\u52A9\u6559"
                        )
                    ),
                    React.createElement(
                        Button,
                        { type: "submit", bsSize: "large", block: true },
                        "\u767B\u5F55"
                    )
                ),
                React.createElement(
                    "form",
                    { action: "/forget/password/id", method: "get" },
                    React.createElement("br", null),
                    React.createElement(
                        Button,
                        { type: "submit", bsSize: "medium", block: true },
                        "\u5FD8\u8BB0\u5BC6\u7801"
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
            Well,
            null,
            React.createElement(
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
            )
        );
    }
});

window.Introduction = React.createClass({
    displayName: "Introduction",

    render: function () {
        return React.createElement(
            "div",
            null,
            React.createElement(
                PageHeader,
                null,
                "\u8F6F\u4EF6\u5DE5\u7A0B"
            ),
            React.createElement(
                "h3",
                null,
                "\xA0\xA0\xA0\xA0\xA0\xA0\xA0\xA0\u6D59\u6C5F\u5927\u5B66\u8F6F\u4EF6\u5B66\u9662\u524D\u8EAB\u662F\u6D59\u6C5F\u5927\u5B66\u8F6F\u4EF6\u4E0E\u7F51\u7EDC\u5B66\u9662\uFF0C\u4E8E2001\u5E742\u670827\u65E5\u5728\u676D\u5DDE\u4E0E\u5B81\u6CE2\u4E24\u5730\u540C\u65F6\u6302\u724C\u6210\u7ACB\uFF0C2001\u5E7412\u6708\u6210\u4E3A\u56FD\u5BB6\u6559\u80B2\u90E8\u548C\u56FD\u5BB6\u53D1\u5C55\u8BA1\u5212\u59D4\u5458\u4F1A\u6279\u51C6\u7684\u9996\u627935\u6240\u56FD\u5BB6\u793A\u8303\u6027\u8F6F\u4EF6\u5B66\u9662\u4E4B\u4E00\uFF0C\u540C\u65F6\u66F4\u540D\u4E3A\u6D59\u6C5F\u5927\u5B66\u8F6F\u4EF6\u5B66\u9662\u3002",
                React.createElement("br", null),
                "\xA0\xA0\xA0\xA0\xA0\xA0\xA0\xA0\u6D59\u6C5F\u5927\u5B66\u56FD\u5BB6\u793A\u8303\u6027\u8F6F\u4EF6\u5B66\u9662\u5206\u522B\u5728\u676D\u5DDE\u548C\u5B81\u6CE2\u529E\u5B66\u3002\u676D\u5DDE\u529E\u5B66\u5730\u70B9\u5728\u6D59\u6C5F\u5927\u5B66\u7389\u6CC9\u6821\u533A\uFF0C\u4EE5\u57F9\u517B\u672C\u79D1\u751F\u4E3A\u4E3B\u3002\u5B81\u6CE2\u529E\u5B66\u5730\u70B9\u5728\u5B81\u6CE2\u56FD\u5BB6\u9AD8\u65B0\u533A\uFF0C\u4EE5\u57F9\u517B\u7814\u7A76\u751F\u4E3A\u4E3B\u3002",
                React.createElement("br", null),
                React.createElement("br", null)
            )
        );
    }
});

window.PanelGroup = ReactBootstrap.PanelGroup;
window.Panel = ReactBootstrap.Panel;
class Friend extends React.Component {
    render() {
        return React.createElement(
            PanelGroup,
            { defaultActiveKey: "1", accordion: true },
            React.createElement(
                Panel,
                { header: "\u53CB\u60C5\u94FE\u63A5", eventKey: "1" },
                this.props.links
            )
        );
    }
}

window.friendLinks = React.createElement(
    "table",
    null,
    React.createElement(
        "tr",
        null,
        React.createElement(
            "td",
            null,
            React.createElement(
                "a",
                { href: "http://zupo.zju.edu.cn", target: "_blank" },
                React.createElement("img", { src: "http://jwbinfosys.zju.edu.cn/images/zupologo.gif", border: "0" })
            )
        ),
        React.createElement("td", { width: "10px" }),
        React.createElement(
            "td",
            null,
            React.createElement(
                "a",
                { href: "http://jwb.zju.edu.cn", target: "_blank" },
                React.createElement("img", { src: "http://jwbinfosys.zju.edu.cn/images/jwblogo.gif", border: "0" })
            )
        ),
        React.createElement("td", { width: "10px" }),
        React.createElement(
            "td",
            null,
            React.createElement(
                "a",
                { href: "http://www.cc98.org/", target: "_blank" },
                React.createElement("img", { src: "http://jwbinfosys.zju.edu.cn/images/cc98logo.gif", border: "0" })
            )
        )
    )
);
