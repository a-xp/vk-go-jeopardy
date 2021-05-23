import {
    Button,
    Cell,
    CellButton,
    Div,
    FormItem,
    FormStatus,
    Group,
    Header,
    IconButton,
    Input,
    Panel,
    PanelHeader,
    PanelHeaderButton,
    PanelHeaderClose,
    Switch
} from "@vkontakte/vkui";
import {
    Icon28AddCircleFillBlue,
    Icon28EditCircleFillBlue,
    Icon28EditOutline,
    Icon28SettingsOutline
} from "@vkontakte/icons";
import {useCallback, useRef} from "react";

function GameEdit({onClose, game, onUpdate, onSelectTopic, onMessageEdit, errors, onSave}) {

    const topicName = useRef();

    const setField = useCallback((k, v) => {
        onUpdate(old => ({...old, [k]: v}))
    }, [onUpdate]);

    const onAddTopic = useCallback(() => {
        onUpdate(old => ({...old, topics: [...old.topics, {name: topicName.current.value, points: 1, q: []}]}));
        topicName.current.value = '';
    }, [onUpdate, topicName]);

    const onRemoveTopic = useCallback((removeIndex) => {
        onUpdate(old => ({...old, topics: old.topics.filter((v, i) => removeIndex !== i)}))
    }, [onUpdate]);

    return (
        <Panel id="game">
            <PanelHeader left={<PanelHeaderClose onClick={onClose}/>} right={<PanelHeaderButton><Icon28SettingsOutline/></PanelHeaderButton>}>{game.name}</PanelHeader>
            <Group>
                {errors && errors.length && <FormItem>
                    <FormStatus header="Некорректно заполнены поля" mode="error">
                        <pre>
                        {errors.join("\n")}
                        </pre>
                    </FormStatus>
                </FormItem>}
                <Group mode="plain">
                    <FormItem top="Название">
                        <Input value={game.name} onChange={e => setField('name', e.currentTarget.value)}/>
                    </FormItem>
                    <FormItem top="Ссылка на стартовый пост">
                        <Input value={game.postLink} onChange={e => setField('postLink', e.currentTarget.value)}/>
                    </FormItem>
                    <FormItem top="Лимит попыток">
                        <Input type="number" value={(game.rules && game.rules.numTries) || ''}
                               onChange={e => setField('numTries', parseInt(e.currentTarget.value))}/>
                    </FormItem>
                    <Cell
                        after={<Switch defaultChecked={game.active} onChange={() => setField('active', !game.active)}/>}
                        description="После запуска игры невозможны изменения">
                        Викторина запущена
                    </Cell>
                    <CellButton before={<Icon28EditOutline/>} onClick={onMessageEdit}>Сообщения бота</CellButton>
                </Group>
                <Group header={<Header mode="tertiary">Темы</Header>} mode="plain">
                    {game.topics && game.topics.map((topic, i) =>
                        <Cell key={i} removable={game.new}
                              onRemove={() => onRemoveTopic(i)}
                              after={<IconButton
                                  onClick={() => onSelectTopic(i + 1)}><Icon28EditCircleFillBlue/></IconButton>}
                              indicator={`${(topic.q && topic.q.length) || 0} вопр.`}>{topic.name}</Cell>
                    )}
                    {game.new && <FormItem>
                        <Input getRef={topicName}
                               after={<IconButton onClick={onAddTopic}><Icon28AddCircleFillBlue/></IconButton>}/>
                    </FormItem>}
                </Group>
                <Div>
                    <Button stretched onClick={onSave}>Сохранить изменения</Button>
                </Div>
            </Group>
        </Panel>
    )
}

export default GameEdit;