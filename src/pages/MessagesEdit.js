import {FormItem, Group, Input, Panel, PanelHeader, PanelHeaderBack} from "@vkontakte/vkui";
import {MESSAGE_CODES} from "../service/default";


function MessagesEdit({messages, onClose, onMessageUpdate}) {
    return (
        <Panel id="messages">
            <PanelHeader left={<PanelHeaderBack onClick={onClose}/>}>Сообщения ботов</PanelHeader>
            <Group>
                {
                    Object.entries(MESSAGE_CODES).map(([key, msg]) => <FormItem key={key} top={msg}>
                        <Input value={messages[key]} onChange={e => onMessageUpdate(key, e.currentTarget.value)}/>
                    </FormItem>)
                }
            </Group>
        </Panel>
    )
}

export default MessagesEdit;