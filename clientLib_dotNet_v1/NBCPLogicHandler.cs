using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace nbcpClientLibv1
{
    public interface NBCPLogicHandler
    {
        public abstract bool WillReAttach();
        public abstract RPCResult RPCHandler(string action, Dictionary<string, object> param);
    }
}
