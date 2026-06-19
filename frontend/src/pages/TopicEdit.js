import {
    Button,
    Div,
    FormItem,
    Group,
    Header,
    IconButton,
    InfoRow,
    Input,
    Panel,
    PanelHeader,
    PanelHeaderBack,
    SimpleCell
} from "@vkontakte/vkui";
import {createRef, useCallback, useEffect, useState} from "react";
import {Icon28AddCircleFillBlue} from "@vkontakte/icons";

function TopicEdit({topic, onEditDone, isNew}) {

    const [curValue, setCurValue] = useState(topic);
    const [ansRefs, setAnsRefs] = useState([]);

    const onDone = useCallback(() => {
        onEditDone(curValue)
    }, [onEditDone, curValue]);

    const onUpdate = useCallback(patch => {
        setCurValue(v => ({...v, ...patch}))
    }, [setCurValue]);

    const onUpdateQuestion = useCallback((qv, qi) => {
        setCurValue(v => ({...v, q: v.q.map((e, i) => qi === i ? {...e, text: qv} : e)}))
    }, [setCurValue]);

    const onAddQuestion = useCallback(() => {
        setCurValue(v => ({...v, q: [...v.q, {text: '', ans: ''}]}))
    }, [setCurValue]);

    const onRemoveQuestion = useCallback((qi) => {
        setCurValue(v => ({...v, q: v.q.filter((e, i) => i !== qi)}))
    }, [setCurValue]);

    const onUpdateAnswer = useCallback((av, qi, ai) => {
        setCurValue(old => ({
            ...old,
            q: old.q.map((q, i) => i === qi ? {...q, ans: q.ans.map((a, j) => j === ai ? av : a)} : q)
        }))
    }, [setCurValue]);

    const onAddAnswer = useCallback((e) => {
        const qi = parseInt(e.currentTarget.getAttribute('data-index'));
        const el = ansRefs[qi].current;
        setCurValue(old => ({...old, q: old.q.map((q, i) => i === qi ? {...q, ans: [...q.ans, el.value]} : q)}))
        el.value = '';
    }, [setCurValue, ansRefs]);

    const onRemoveAnswer = useCallback((qi, ai) => {
        setCurValue(old => ({
            ...old,
            q: old.q.map((q, i) => i === qi ? {...q, ans: q.ans.filter((a, j) => j !== ai)} : q)
        }))
    }, [setCurValue]);

    const numQ = (curValue.q && curValue.q.length) || 0;

    useEffect(() => {
        setAnsRefs(prev => Array(numQ).fill().map((_, i) => prev[i] || createRef()));
    }, [numQ, setAnsRefs])

    if (!topic) {
        return null;
    }

    return (
        <Panel id="topic">
            <PanelHeader left={<PanelHeaderBack onClick={onDone}/>}>{topic.name}</PanelHeader>
            <Group>
                <Group mode="plain">
                    <FormItem top="Название темы">
                        <Input required value={curValue.name} onChange={e => onUpdate({name: e.currentTarget.value})}/>
                    </FormItem>
                    {isNew ? <FormItem top="Баллы">
                            <Input value={curValue.points} type="number"
                                   onChange={e => onUpdate({points: parseInt(e.currentTarget.value) || 0})}/>
                        </FormItem> :
                        <SimpleCell>
                            <InfoRow header="Баллы">
                                {curValue.points}
                            </InfoRow>
                        </SimpleCell>}
                </Group>
                {curValue.q && curValue.q.map((item, i) => <Group mode="plain" key={i} header={<Header
                    mode="tertiary">Вопрос {i + 1}</Header>}>
                    <FormItem removable={isNew} onRemove={() => onRemoveQuestion(i)}>
                        <Input required value={item.text} onChange={e => onUpdateQuestion(e.currentTarget.value, i)}/>
                    </FormItem>
                    {item.ans && item.ans.map((ans, j) => <FormItem key={j} removable
                                                                    onRemove={() => onRemoveAnswer(i, j)}>
                        <Input required value={ans} onChange={e => onUpdateAnswer(e.currentTarget.value, i, j)}/>
                    </FormItem>)}
                    <FormItem>
                        <Input getRef={ansRefs[i]} after={<IconButton data-index={i}
                                                                      onClick={onAddAnswer}><Icon28AddCircleFillBlue/></IconButton>}/>
                    </FormItem>
                </Group>)}
                {isNew && <Div>
                    <Button stretched onClick={onAddQuestion}>Добавить вопрос</Button>
                </Div>}
            </Group>
        </Panel>
    )
}

export default TopicEdit;