import { FullScreenLoading } from "@/components/common";
import { useEffect } from "react";
import { toast } from "sonner";

export const CallbackPage = () => {
  const params = new URLSearchParams(window.location.search);
  const toastmsg = params.get("toast");
  const redirect = params.get("redirect") || "/";

  useEffect(() => {
    if (toastmsg) {
      toast(toastmsg);
      setTimeout(() => {
        window.location.href = redirect;
      }, 1500);
    } else {
      window.location.href = redirect;
    }
  }, []);

  return <FullScreenLoading />;
};
