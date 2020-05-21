## COVID-19 India (Server - Golang)

Refined Stats From COVID-19 Official Website

`Data is periodontally collected from Official Website`

#### API Server:

`https://api.novelcoronaindia.info`

#### API Endpoints:

1. `GET /covid19/all`
2. `GET /covid19/date/{date}`
3. `GET /covid19/dateRange/{startDate}/{endDate}`
4. `GET /covid19/formattedData` (Optional query parameter: `startDate` and `endDate`).

   - `/covid19/formattedData?startDate=15-05-2020` (Optional endDate. Default: `Today`)
   - `/covid19/formattedData?endDate=21-05-2020` (Optional startDate. Default: `03-04-2020`)
   - `/covid19/formattedData?startDate=15-05-2020&endDate=21-05-2020`
   - Sample Response

   ```js
   {
     "data": [
       {
         "date": "21-05-2020",
         "confirmed": 1000,
         "recovered": 500,
         "death": 5,
         "active": 499,
         "stateWise": [
           {
             "stateName": "West Bengal",
             "confirmed": 500,
             "recovered": 200,
             "death": 1,
             "active": 299,
           },
           // more state data
         ],
       },
    // more date wise data
     ];
   }
   ```

Note:

1. Accepted date format: `DD-MM-YYYY` (eg. `18-04-2020`)
2. Data available from `03-04-2020`

Frontend Code: https://github.com/ritwickdey/covid-19-india-react

[MIT Licence](./LICENCE)
