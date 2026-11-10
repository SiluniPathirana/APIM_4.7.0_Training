package org.pizzashack.handler;

import java.util.Map;

import org.apache.synapse.MessageContext;
import org.apache.synapse.core.axis2.Axis2MessageContext;
import org.apache.synapse.core.axis2.Axis2Sender;
import org.apache.synapse.rest.AbstractHandler;

/**
 * Sample API Gateway handler for the training lab.
 * Lets a request through only when its "Username" header starts with "admin".
 * Any other request is answered with HTTP 401 and is not sent to the back end.
 */
public class UsernameHeaderHandler extends AbstractHandler {

    private static final String HEADER = "Username";

    @Override
    public boolean handleRequest(MessageContext messageContext) {
        org.apache.axis2.context.MessageContext axis2Ctx =
                ((Axis2MessageContext) messageContext).getAxis2MessageContext();
        Map<?, ?> headers = (Map<?, ?>) axis2Ctx.getProperty(
                org.apache.axis2.context.MessageContext.TRANSPORT_HEADERS);
        Object user = headers == null ? null : headers.get(HEADER);
        System.out.println("UsernameHeaderHandler: " + HEADER + " = " + user);
        if (user != null && user.toString().startsWith("admin")) {
            return true;
        }
        // Reject: send a 401 with no body back to the client instead of calling the back end.
        axis2Ctx.setProperty("HTTP_SC", 401);
        axis2Ctx.setProperty("NO_ENTITY_BODY", Boolean.TRUE);
        messageContext.setProperty("RESPONSE", "true");
        messageContext.setTo(null);
        Axis2Sender.sendBack(messageContext);
        return false;
    }

    @Override
    public boolean handleResponse(MessageContext messageContext) {
        return true;
    }
}
