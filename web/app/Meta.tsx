"use client";

import { useEffect, useState } from "react";
import Providers from "./providers";
import Button from "./MetaButton";

const Meta = ({ dict }: { dict: any }) => {
  const [href, setHref] = useState("");
  useEffect(() => {
    setHref(window.location.href);
  }, []);

  return (
    <Providers href={href}>
      <Button title={dict.title} />
    </Providers>
  );
};

export default Meta;
