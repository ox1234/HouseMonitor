package template

var HouseTemplate = `
{
  "config": {
    "wide_screen_mode": true
  },
  "elements": [
    {
      "tag": "div",
      "text": {
        "content": "%s",
        "tag": "lark_md"
      }
    }
  ],
  "header": {
    "template": "blue",
    "title": {
      "content": "🏘️ 豆瓣租房小组房源上新了！",
      "tag": "plain_text"
    }
  }
}
`
