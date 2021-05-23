import React, {useCallback, useState} from "react";
import GameEdit from "./GameEdit";
import TopicEdit from "./TopicEdit";
import MessagesEdit from "./MessagesEdit";
import {formatPostId, hasText, parsePostId} from "../service/util";
import {View} from "@vkontakte/vkui";
import {MESSAGE_CODES} from "../service/default";

function validatePayload(newPayload) {
    const errors = [];
    console.log(newPayload);
    if (newPayload.active) {
        if (!hasText(newPayload.name)) {
            errors.push('Название не заполнено');
        }
        if (!newPayload.post) {
            errors.push('Неправильна ссылка на стартовый пост')
        }
        if (!newPayload.topics || !newPayload.topics.length) {
            errors.push('Не заданы темы вопросов')
        } else {
            newPayload.topics.forEach((topic, i) => {
                if (!hasText(topic.name)) {
                    errors.push(`Не задано название темы ${i + 1}`)
                }
                if (!topic.points || topic.points < 1) {
                    errors.push(`Не заданы баллы темы ${i + 1}`)
                }
                if (!topic.q || !topic.q.length) {
                    errors.push(`Не заданы вопросы темы ${i + 1}`)
                } else {
                    topic.q.forEach((q, j) => {
                        if (!hasText(q.text)) {
                            errors.push(`Не задан вопрос ${j + 1} в теме ${i + 1}`)
                        }
                        if (!q.ans || !q.ans.length) {
                            errors.push(`Не задан ответ ${j + 1} в теме ${i + 1}`)
                        }
                    });
                }
            });
        }
        Object.keys(MESSAGE_CODES).forEach(k => {
            if (!hasText(newPayload.messages[k])) {
                errors.push('Не задано сообщение бота: ' + MESSAGE_CODES[k])
            }
        })

    }
    return errors;
}


function GameEditWrapper({game, onClose, onSave}) {

    const [newValue, setNewValue] = useState(() => (
        {...game, postLink: formatPostId(game.post), numTries: (game.rules && game.rules.numTries) || 0}
    ));

    const [editMessage, setEditMessage] = useState(false);
    const [showTopic, setShowTopic] = useState();
    const [errors, setErrors] = useState();

    const onSelectTopic = useCallback((i) => {
        setShowTopic(i)
    }, [setShowTopic]);

    const onTopicUpdate = useCallback((newValue) => {
        setNewValue(v => ({...v, topics: v.topics.map((e, i) => i === showTopic - 1 ? newValue : e)}))
        setShowTopic(0)
    }, [setNewValue, showTopic]);

    const onMessageEdit = useCallback(() => {
        setEditMessage(true);
    }, [setEditMessage]);

    const onMessageEditDone = useCallback(() => {
        setEditMessage(false);
    }, [setEditMessage]);

    const onMessageUpdate = useCallback((key, value) => {
        setNewValue(old => ({...old, messages: {...old.messages, [key]: value}}))
    }, [setNewValue]);

    const onSaveInt = useCallback(() => {
        setErrors(null);
        const value = {
            ...newValue,
            post: parsePostId(newValue.postLink),
            "new": newValue.new && !newValue.active,
            rules: {numTries: newValue.numTries}
        };
        const err = validatePayload(value);
        if (err.length) {
            setErrors(err);
        } else {
            onSave(value);
        }
    }, [onSave, newValue, setErrors]);

    return (
        <View activePanel={editMessage ? "messages" : (showTopic ? "topic" : "main")}>
            <MessagesEdit id="messages" messages={newValue.messages} onClose={onMessageEditDone}
                          onMessageUpdate={onMessageUpdate}/>
            <GameEdit id="main" errors={errors} onSave={onSaveInt} game={newValue} onUpdate={setNewValue}
                      onClose={onClose} onSelectTopic={onSelectTopic} onMessageEdit={onMessageEdit}/>
            <TopicEdit id="topic" topic={showTopic && newValue.topics[showTopic - 1]} onEditDone={onTopicUpdate}
                       isNew={game.new}/>
        </View>
    )
}

export default GameEditWrapper;

