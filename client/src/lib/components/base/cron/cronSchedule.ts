export type CronSchedule =
  | { mode: "minutes"; every: number }
  | { mode: "hourly"; minute: number }
  | { mode: "daily"; hour: number; minute: number }
  | { mode: "weekly"; weekday: number; hour: number; minute: number }
  | { mode: "monthly"; day: number; hour: number; minute: number }
  | { mode: "custom"; expression: string };

export const weekdayNames = [
  "Sunday",
  "Monday",
  "Tuesday",
  "Wednesday",
  "Thursday",
  "Friday",
  "Saturday",
];

export const scheduleModes: { value: string; label: string }[] = [
  { value: "minutes", label: "Every few minutes" },
  { value: "hourly", label: "Every hour" },
  { value: "daily", label: "Every day" },
  { value: "weekly", label: "Every week" },
  { value: "monthly", label: "Every month" },
  { value: "custom", label: "Custom cron expression" },
];

const pad = (value: number) => String(value).padStart(2, "0");
const time = (hour: number, minute: number) => `${pad(hour)}:${pad(minute)}`;
const asInt = (value: string) => (/^\d+$/.test(value) ? Number(value) : null);
const clamp = (value: number, min: number, max: number) =>
  Math.min(max, Math.max(min, Math.round(value || min)));

export function toCronExpression(schedule: CronSchedule): string {
  switch (schedule.mode) {
    case "minutes":
      return `*/${clamp(schedule.every, 1, 59)} * * * *`;
    case "hourly":
      return `${clamp(schedule.minute, 0, 59)} * * * *`;
    case "daily":
      return `${clamp(schedule.minute, 0, 59)} ${clamp(schedule.hour, 0, 23)} * * *`;
    case "weekly":
      return `${clamp(schedule.minute, 0, 59)} ${clamp(schedule.hour, 0, 23)} * * ${clamp(schedule.weekday, 0, 6)}`;
    case "monthly":
      return `${clamp(schedule.minute, 0, 59)} ${clamp(schedule.hour, 0, 23)} ${clamp(schedule.day, 1, 31)} * *`;
    case "custom":
      return schedule.expression.trim();
  }
}

export function describeSchedule(schedule: CronSchedule): string {
  switch (schedule.mode) {
    case "minutes":
      return schedule.every === 1
        ? "Every minute"
        : `Every ${schedule.every} minutes`;
    case "hourly":
      return `Every hour at minute ${pad(schedule.minute)}`;
    case "daily":
      return `Every day at ${time(schedule.hour, schedule.minute)}`;
    case "weekly":
      return `Every ${weekdayNames[schedule.weekday]} at ${time(schedule.hour, schedule.minute)}`;
    case "monthly":
      return `Day ${schedule.day} of every month at ${time(schedule.hour, schedule.minute)}`;
    case "custom":
      return "Custom cron expression";
  }
}

export function isCronExpression(value: string): boolean {
  const fields = value.trim().split(/\s+/);
  if (fields.length !== 5 && fields.length !== 6) return false;
  return fields.every((field) => /^[\d*/,\-]+$/.test(field));
}

export function describeCronExpression(expression: string): string {
  const schedule = parseCronExpression(expression);
  return schedule ? describeSchedule(schedule) : "Custom cron expression";
}

export function parseCronExpression(expression: string): CronSchedule | null {
  const fields = expression.trim().split(/\s+/);
  if (fields.length !== 5) return null;

  const [minute, hour, day, month, weekday] = fields;

  if (
    /^\*\/\d+$/.test(minute) &&
    hour === "*" &&
    day === "*" &&
    month === "*" &&
    weekday === "*"
  ) {
    return { mode: "minutes", every: Number(minute.slice(2)) };
  }

  const minuteValue = asInt(minute);
  if (minuteValue === null) return { mode: "custom", expression };

  if (hour === "*" && day === "*" && month === "*" && weekday === "*") {
    return { mode: "hourly", minute: minuteValue };
  }

  const hourValue = asInt(hour);
  if (hourValue === null) return { mode: "custom", expression };

  if (day === "*" && month === "*" && weekday === "*") {
    return { mode: "daily", hour: hourValue, minute: minuteValue };
  }

  if (day === "*" && month === "*") {
    const weekdayValue = asInt(weekday);
    if (weekdayValue !== null) {
      return {
        mode: "weekly",
        weekday: weekdayValue,
        hour: hourValue,
        minute: minuteValue,
      };
    }
  }

  if (month === "*" && weekday === "*") {
    const dayValue = asInt(day);
    if (dayValue !== null) {
      return {
        mode: "monthly",
        day: dayValue,
        hour: hourValue,
        minute: minuteValue,
      };
    }
  }

  return { mode: "custom", expression };
}
