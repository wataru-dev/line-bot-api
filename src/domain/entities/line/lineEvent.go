package line

type LineEvent struct {
    Type       string `json:"type"`
    ReplyToken string `json:"replyToken"`
    Source     struct {
        UserID string `json:"userId"`
    } `json:"source"`
    Message struct {
        ID   string `json:"id"`
        Type string `json:"type"`
        Text string `json:"text"`
    } `json:"message"`
}

type LineWebhook struct {
    Events []LineEvent `json:"events"`
}

type ReplyMessage struct {
    ReplyToken string    `json:"replyToken"`
    Messages   []Message `json:"messages"`
}

type Message struct {
    Type string `json:"type"`
    Text string `json:"text"`
}
