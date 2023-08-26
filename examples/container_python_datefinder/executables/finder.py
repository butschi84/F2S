import datefinder
import sys
import json
from datetime import datetime, timezone

# get input args
arguments = sys.argv[1:]
inputText = arguments[0]

# use datefinder to find all dates
matches = datefinder.find_dates(inputText, source=True)


matched_dates = []
for match in matches:
    matched_dates.append({
        "date": match[0],
        "date_timestamp": int(match[0].astimezone(timezone.utc).timestamp()),
        "raw_date": match[1],
        "raw_offset": inputText.find(match[1]),
        "raw_length": len(match[1])
    })

json_output = json.dumps(matched_dates, default=str)

# Print the JSON array
print(json_output)
