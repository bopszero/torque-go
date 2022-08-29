package docs

var doc = "" +
	"openapi: 3.0.0\n" +
	"info:\n" +
	"  version: 1.0.0\n" +
	"  title: Torque Transaction Service\n" +
	"servers:\n" +
	"  - url: /\n" +
	"components:\n" +
	"  schemas:\n" +
	"    ResponseCode:\n" +
	"      type: string\n" +
	"      example: success\n" +
	"    ResponseMessage:\n" +
	"      type: string\n" +
	"      example: \"Field validation for 'user_id' failed on the 'required' tag.\"\n" +
	"    Decimal:\n" +
	"      type: string\n" +
	"      description: Decimal number in string\n" +
	"      example: \"14211.364041\"\n" +
	"    ID:\n" +
	"      type: integer\n" +
	"      minimum: 1\n" +
	"      example: 8466032601\n" +
	"    Date:\n" +
	"      type: string\n" +
	"      format: date\n" +
	"      example: \"2020-12-01\"\n" +
	"    Reference:\n" +
	"      type: string\n" +
	"      example: \"123456789\"\n" +
	"    Currency:\n" +
	"      type: string\n" +
	"      enum:\n" +
	"        - BTC\n" +
	"        - BCH\n" +
	"        - LTC\n" +
	"        - ETH\n" +
	"        - TRX\n" +
	"        - USDT\n" +
	"        - TORQ\n" +
	"    Timestamp:\n" +
	"      type: integer\n" +
	"      example: 1578396200\n" +
	"    Note:\n" +
	"      type: string\n" +
	"      example: I lay my love on you...\n" +
	"    Network:\n" +
	"      type: string\n" +
	"      enum:\n" +
	"        - BTC\n" +
	"        - BTC:TESTNET\n" +
	"        - BCH\n" +
	"        - BCH:TESTNET\n" +
	"        - LTC\n" +
	"        - LTC:TESTNET\n" +
	"        - ETH\n" +
	"        - ETH:TEST_ROPSTEN\n" +
	"        - TRX\n" +
	"        - TRX:TEST_SHASTA\n" +
	"        - XRP\n" +
	"        - XRP:TESTNET\n" +
	"        - TORQ\n" +
	"    UserTierType:\n" +
	"      type: integer\n" +
	"      enum:\n" +
	"        - 1\n" +
	"        - 2\n" +
	"        - 3\n" +
	"        - 4\n" +
	"        - 6\n" +
	"        - 7\n" +
	"      description: |\n" +
	"        * 1 - Agent\n" +
	"        * 2 - Senior Partner / Market Leader\n" +
	"        * 3 - Regional Leader\n" +
	"        * 4 - Governor\n" +
	"        * 6 - Global Leader\n" +
	"        * 7 - Mentor\n" +
	"    TypeCode:\n" +
	"      type: string\n" +
	"      enum:\n" +
	"        - cr_system\n" +
	"        - cr_transfer\n" +
	"        - cr_deposit\n" +
	"        - cr_withdraw_investment_reverse\n" +
	"        - cr_daily_profit\n" +
	"        - cr_affiliate_commission\n" +
	"        - cr_leader_commission\n" +
	"        - cr_reinvest_src_reverse\n" +
	"        - cr_reinvest_dst\n" +
	"        - cr_withdraw_profit_reverse\n" +
	"        - cr_promo_code\n" +
	"        - cr_product_travel_ticket_refund\n" +
	"        - cr_product_mall_order_refund\n" +
	"        - cr_product_flight_ticket_refund\n" +
	"        - dr_system\n" +
	"        - dr_transfer\n" +
	"        - dr_deposit_reverse\n" +
	"        - dr_withdraw_investment\n" +
	"        - dr_daily_profit_reverse\n" +
	"        - dr_affiliate_commission_reverse\n" +
	"        - dr_leader_commission_reverse\n" +
	"        - dr_reinvest_src\n" +
	"        - dr_withdraw_profit\n" +
	"        - dr_product_event_ticket\n" +
	"        - dr_product_travel_ticket\n" +
	"        - dr_product_mall_order\n" +
	"        - dr_product_flight_ticket\n" +
	"    Address:\n" +
	"      type: string\n" +
	"      minLength: 16\n" +
	"      maxLength: 255\n" +
	"      example: bc1q74m8sy7dpqwuegcrddfwqp3gzwau62k7rys0my\n" +
	"    TxnHash:\n" +
	"      type: string\n" +
	"      minLength: 32\n" +
	"      maxLength: 128\n" +
	"      example: 947cfb43f4d3bea439f5045f9c47e811a47dae26e2976be0ef56d5e058f81e0a\n" +
	"    FeeInfo:\n" +
	"      type: object\n" +
	"      properties:\n" +
	"        currency:\n" +
	"          $ref: \"#/components/schemas/Currency\"\n" +
	"        price:\n" +
	"          $ref: \"#/components/schemas/Decimal\"\n" +
	"        price_high:\n" +
	"          $ref: \"#/components/schemas/Decimal\"\n" +
	"        price_low:\n" +
	"          $ref: \"#/components/schemas/Decimal\"\n" +
	"        limit_max_quantity:\n" +
	"          type: integer\n" +
	"        limit_min_value:\n" +
	"          $ref: \"#/components/schemas/Decimal\"\n" +
	"        limit_max_value:\n" +
	"          $ref: \"#/components/schemas/Decimal\"\n" +
	"        base_currency:\n" +
	"          $ref: \"#/components/schemas/Currency\"\n" +
	"        to_base_multiple:\n" +
	"          $ref: \"#/components/schemas/Decimal\"\n" +
	"      required:\n" +
	"        - currency\n" +
	"        - price\n" +
	"        - price_high\n" +
	"        - price_low\n" +
	"        - limit_max_quantity\n" +
	"        - limit_min_value\n" +
	"        - limit_max_value\n" +
	"    BonusPoolLeaderTierInfo:\n" +
	"      type: object\n" +
	"      properties:\n" +
	"        tier_type:\n" +
	"          $ref: \"#/components/schemas/UserTierType\"\n" +
	"        uids:\n" +
	"          type: array\n" +
	"          items:\n" +
	"            $ref: \"#/components/schemas/ID\"\n" +
	"        rate:\n" +
	"          $ref: \"#/components/schemas/Decimal\"\n" +
	"        additional_user_info_list:\n" +
	"          type: array\n" +
	"          items:\n" +
	"            type: object\n" +
	"            properties:\n" +
	"              uid:\n" +
	"                $ref: \"#/components/schemas/ID\"\n" +
	"              note:\n" +
	"                $ref: \"#/components/schemas/Note\"\n" +
	"      required:\n" +
	"        - tier_type\n" +
	"        - uids\n" +
	"        - rate\n" +
	"    BonusPoolLeaderExecutionStatus:\n" +
	"      type: integer\n" +
	"      enum:\n" +
	"        - 1\n" +
	"        - 2\n" +
	"        - 10\n" +
	"      description: |\n" +
	"        * 1 - Init\n" +
	"        * 2 - Executing\n" +
	"        * 10 - Completed\n" +
	"paths:\n" +
	"  /v1/balance/user/get/:\n" +
	"    post:\n" +
	"      summary: Get a User balance\n" +
	"      tags:\n" +
	"        - Balance\n" +
	"      requestBody:\n" +
	"        required: true\n" +
	"        content:\n" +
	"          application/json:\n" +
	"            schema:\n" +
	"              type: object\n" +
	"              properties:\n" +
	"                currency:\n" +
	"                  $ref: \"#/components/schemas/Currency\"\n" +
	"                uid:\n" +
	"                  $ref: \"#/components/schemas/ID\"\n" +
	"              required:\n" +
	"                - uid\n" +
	"      responses:\n" +
	"        \"200\":\n" +
	"          description: OK\n" +
	"          content:\n" +
	"            application/json:\n" +
	"              schema:\n" +
	"                type: object\n" +
	"                properties:\n" +
	"                  message:\n" +
	"                    $ref: \"#/components/schemas/ResponseMessage\"\n" +
	"                  data:\n" +
	"                    type: object\n" +
	"                    properties:\n" +
	"                      balances:\n" +
	"                        type: array\n" +
	"                        items:\n" +
	"                          type: object\n" +
	"                          properties:\n" +
	"                            uid:\n" +
	"                              $ref: \"#/components/schemas/ID\"\n" +
	"                            currency:\n" +
	"                              $ref: \"#/components/schemas/Currency\"\n" +
	"                            amount:\n" +
	"                              $ref: \"#/components/schemas/Decimal\"\n" +
	"                            update_time:\n" +
	"                              $ref: \"#/components/schemas/Timestamp\"\n" +
	"  /v1/balance/txn/add/:\n" +
	"    post:\n" +
	"      summary: Add a Balance Transaction\n" +
	"      tags:\n" +
	"        - Balance\n" +
	"      requestBody:\n" +
	"        required: true\n" +
	"        content:\n" +
	"          application/json:\n" +
	"            schema:\n" +
	"              type: object\n" +
	"              properties:\n" +
	"                currency:\n" +
	"                  $ref: \"#/components/schemas/Currency\"\n" +
	"                uid:\n" +
	"                  $ref: \"#/components/schemas/ID\"\n" +
	"                amount:\n" +
	"                  $ref: \"#/components/schemas/Decimal\"\n" +
	"                type_code:\n" +
	"                  $ref: \"#/components/schemas/TypeCode\"\n" +
	"                ref:\n" +
	"                  $ref: \"#/components/schemas/Reference\"\n" +
	"              required:\n" +
	"                - currency\n" +
	"                - uid\n" +
	"                - amount\n" +
	"                - type\n" +
	"                - ref\n" +
	"      responses:\n" +
	"        \"200\":\n" +
	"          description: OK\n" +
	"          content:\n" +
	"            application/json:\n" +
	"              schema:\n" +
	"                type: object\n" +
	"                properties:\n" +
	"                  message:\n" +
	"                    $ref: \"#/components/schemas/ResponseMessage\"\n" +
	"                  data:\n" +
	"                    type: object\n" +
	"                    nullable: true\n" +
	"                    properties:\n" +
	"                      txn:\n" +
	"                        type: object\n" +
	"                        properties:\n" +
	"                          id:\n" +
	"                            $ref: \"#/components/schemas/ID\"\n" +
	"                          user_id:\n" +
	"                            $ref: \"#/components/schemas/ID\"\n" +
	"                          parent_id:\n" +
	"                            $ref: \"#/components/schemas/ID\"\n" +
	"                          currency:\n" +
	"                            $ref: \"#/components/schemas/Currency\"\n" +
	"                          amount:\n" +
	"                            $ref: \"#/components/schemas/Decimal\"\n" +
	"                          balance:\n" +
	"                            $ref: \"#/components/schemas/Decimal\"\n" +
	"                          type_code:\n" +
	"                            $ref: \"#/components/schemas/TypeCode\"\n" +
	"                          ref:\n" +
	"                            $ref: \"#/components/schemas/Reference\"\n" +
	"  /v1/payment/investment/deposit/account/get/:\n" +
	"    post:\n" +
	"      summary: Get user Deposit account\n" +
	"      tags:\n" +
	"        - Payment\n" +
	"      requestBody:\n" +
	"        required: true\n" +
	"        content:\n" +
	"          application/json:\n" +
	"            schema:\n" +
	"              type: object\n" +
	"              properties:\n" +
	"                uid:\n" +
	"                  $ref: \"#/components/schemas/ID\"\n" +
	"                currency:\n" +
	"                  $ref: \"#/components/schemas/Currency\"\n" +
	"                network:\n" +
	"                  $ref: \"#/components/schemas/Network\"\n" +
	"              required:\n" +
	"                - uid\n" +
	"                - currency\n" +
	"      responses:\n" +
	"        \"200\":\n" +
	"          description: OK\n" +
	"          content:\n" +
	"            application/json:\n" +
	"              schema:\n" +
	"                type: object\n" +
	"                properties:\n" +
	"                  message:\n" +
	"                    $ref: \"#/components/schemas/ResponseMessage\"\n" +
	"                  data:\n" +
	"                    type: object\n" +
	"                    properties:\n" +
	"                      uid:\n" +
	"                        $ref: \"#/components/schemas/ID\"\n" +
	"                      currency:\n" +
	"                        $ref: \"#/components/schemas/Currency\"\n" +
	"                      network:\n" +
	"                        $ref: \"#/components/schemas/Network\"\n" +
	"                      address:\n" +
	"                        $ref: \"#/components/schemas/Address\"\n" +
	"                      create_time:\n" +
	"                        $ref: \"#/components/schemas/Timestamp\"\n" +
	"  /v1/payment/investment/deposit/crawl/:\n" +
	"    post:\n" +
	"      summary: Crawl deposits from a Block\n" +
	"      tags:\n" +
	"        - Payment\n" +
	"      requestBody:\n" +
	"        required: true\n" +
	"        content:\n" +
	"          application/json:\n" +
	"            schema:\n" +
	"              type: object\n" +
	"              properties:\n" +
	"                currency:\n" +
	"                  $ref: \"#/components/schemas/Currency\"\n" +
	"                network:\n" +
	"                  $ref: \"#/components/schemas/Network\"\n" +
	"                block_height:\n" +
	"                  type: integer\n" +
	"              required:\n" +
	"                - currency\n" +
	"                - network\n" +
	"                - block_height\n" +
	"      responses:\n" +
	"        \"200\":\n" +
	"          description: OK\n" +
	"          content:\n" +
	"            application/json:\n" +
	"              schema:\n" +
	"                type: object\n" +
	"                properties:\n" +
	"                  code:\n" +
	"                    $ref: \"#/components/schemas/ResponseCode\"\n" +
	"                  message:\n" +
	"                    $ref: \"#/components/schemas/ResponseMessage\"\n" +
	"  /v1/payment/investment/deposit/submit/:\n" +
	"    post:\n" +
	"      summary: Submit a new Deposit\n" +
	"      tags:\n" +
	"        - Payment\n" +
	"      requestBody:\n" +
	"        required: true\n" +
	"        content:\n" +
	"          application/json:\n" +
	"            schema:\n" +
	"              type: object\n" +
	"              properties:\n" +
	"                uid:\n" +
	"                  $ref: \"#/components/schemas/ID\"\n" +
	"                currency:\n" +
	"                  $ref: \"#/components/schemas/Currency\"\n" +
	"                network:\n" +
	"                  $ref: \"#/components/schemas/Network\"\n" +
	"                txn_hash:\n" +
	"                  $ref: \"#/components/schemas/TxnHash\"\n" +
	"                txn_index:\n" +
	"                  type: integer\n" +
	"                address:\n" +
	"                  $ref: \"#/components/schemas/Address\"\n" +
	"                amount:\n" +
	"                  $ref: \"#/components/schemas/Decimal\"\n" +
	"              required:\n" +
	"                - uid\n" +
	"                - currency\n" +
	"                - txn_hash\n" +
	"                - address\n" +
	"                - amount\n" +
	"      responses:\n" +
	"        \"200\":\n" +
	"          description: OK\n" +
	"          content:\n" +
	"            application/json:\n" +
	"              schema:\n" +
	"                type: object\n" +
	"                properties:\n" +
	"                  code:\n" +
	"                    $ref: \"#/components/schemas/ResponseCode\"\n" +
	"                  message:\n" +
	"                    $ref: \"#/components/schemas/ResponseMessage\"\n" +
	"                  data:\n" +
	"                    type: object\n" +
	"                    properties:\n" +
	"                      id:\n" +
	"                        $ref: \"#/components/schemas/ID\"\n" +
	"                      coin_id:\n" +
	"                        type: integer\n" +
	"                        example: 2\n" +
	"                      currency:\n" +
	"                        $ref: \"#/components/schemas/Currency\"\n" +
	"                      network:\n" +
	"                        $ref: \"#/components/schemas/Network\"\n" +
	"                      uid:\n" +
	"                        $ref: \"#/components/schemas/ID\"\n" +
	"                      status:\n" +
	"                        type: string\n" +
	"                        example: \"Under Processing\"\n" +
	"                      txn_hash:\n" +
	"                        $ref: \"#/components/schemas/TxnHash\"\n" +
	"                      txn_index:\n" +
	"                        type: integer\n" +
	"                      address:\n" +
	"                        $ref: \"#/components/schemas/Address\"\n" +
	"                      amount:\n" +
	"                        $ref: \"#/components/schemas/Decimal\"\n" +
	"                      create_time:\n" +
	"                        $ref: \"#/components/schemas/Timestamp\"\n" +
	"  /v1/payment/investment/deposit/approve/:\n" +
	"    post:\n" +
	"      summary: Approve a Deposit\n" +
	"      tags:\n" +
	"        - Payment\n" +
	"      requestBody:\n" +
	"        required: true\n" +
	"        content:\n" +
	"          application/json:\n" +
	"            schema:\n" +
	"              type: object\n" +
	"              properties:\n" +
	"                id:\n" +
	"                  $ref: \"#/components/schemas/ID\"\n" +
	"                note:\n" +
	"                  $ref: \"#/components/schemas/Note\"\n" +
	"              required:\n" +
	"                - id\n" +
	"      responses:\n" +
	"        \"200\":\n" +
	"          description: OK\n" +
	"          content:\n" +
	"            application/json:\n" +
	"              schema:\n" +
	"                type: object\n" +
	"                properties:\n" +
	"                  code:\n" +
	"                    $ref: \"#/components/schemas/ResponseCode\"\n" +
	"                  message:\n" +
	"                    $ref: \"#/components/schemas/ResponseMessage\"\n" +
	"  /v1/payment/investment/withdraw/submit/:\n" +
	"    post:\n" +
	"      summary: Add a Investment Withdraw request\n" +
	"      tags:\n" +
	"        - Payment\n" +
	"      requestBody:\n" +
	"        required: true\n" +
	"        content:\n" +
	"          application/json:\n" +
	"            schema:\n" +
	"              type: object\n" +
	"              properties:\n" +
	"                uid:\n" +
	"                  $ref: \"#/components/schemas/ID\"\n" +
	"                amount:\n" +
	"                  $ref: \"#/components/schemas/Decimal\"\n" +
	"                currency:\n" +
	"                  $ref: \"#/components/schemas/Currency\"\n" +
	"                address:\n" +
	"                  $ref: \"#/components/schemas/Address\"\n" +
	"              required:\n" +
	"                - currency\n" +
	"                - uid\n" +
	"                - amount\n" +
	"                - address\n" +
	"      responses:\n" +
	"        \"200\":\n" +
	"          description: OK\n" +
	"          content:\n" +
	"            application/json:\n" +
	"              schema:\n" +
	"                type: object\n" +
	"                properties:\n" +
	"                  message:\n" +
	"                    $ref: \"#/components/schemas/ResponseMessage\"\n" +
	"                  data:\n" +
	"                    type: object\n" +
	"                    properties:\n" +
	"                      withdraw_request:\n" +
	"                        type: object\n" +
	"                        properties:\n" +
	"                          id:\n" +
	"                            $ref: \"#/components/schemas/ID\"\n" +
	"                          code:\n" +
	"                            type: string\n" +
	"                            example: C-200508-939264\n" +
	"                          uid:\n" +
	"                            $ref: \"#/components/schemas/ID\"\n" +
	"                          amount:\n" +
	"                            $ref: \"#/components/schemas/Decimal\"\n" +
	"                          currency:\n" +
	"                            $ref: \"#/components/schemas/Currency\"\n" +
	"                          address:\n" +
	"                            $ref: \"#/components/schemas/Address\"\n" +
	"                          status:\n" +
	"                            type: string\n" +
	"                          create_time:\n" +
	"                            $ref: \"#/components/schemas/Timestamp\"\n" +
	"  /v1/payment/investment/withdraw/reject/:\n" +
	"    post:\n" +
	"      summary: Reject a Investment Withdraw request\n" +
	"      tags:\n" +
	"        - Payment\n" +
	"      requestBody:\n" +
	"        required: true\n" +
	"        content:\n" +
	"          application/json:\n" +
	"            schema:\n" +
	"              type: object\n" +
	"              properties:\n" +
	"                id:\n" +
	"                  $ref: \"#/components/schemas/ID\"\n" +
	"                note:\n" +
	"                  $ref: \"#/components/schemas/Note\"\n" +
	"              required:\n" +
	"                - id\n" +
	"      responses:\n" +
	"        \"200\":\n" +
	"          description: OK\n" +
	"          content:\n" +
	"            application/json:\n" +
	"              schema:\n" +
	"                type: object\n" +
	"                properties:\n" +
	"                  message:\n" +
	"                    $ref: \"#/components/schemas/ResponseMessage\"\n" +
	"                  data:\n" +
	"                    type: object\n" +
	"                    properties:\n" +
	"                      ok:\n" +
	"                        type: boolean\n" +
	"  /v1/payment/profit/withdraw/submit/:\n" +
	"    post:\n" +
	"      summary: Add a Profit Withdraw request\n" +
	"      tags:\n" +
	"        - Payment\n" +
	"      requestBody:\n" +
	"        required: true\n" +
	"        content:\n" +
	"          application/json:\n" +
	"            schema:\n" +
	"              type: object\n" +
	"              properties:\n" +
	"                uid:\n" +
	"                  $ref: \"#/components/schemas/ID\"\n" +
	"                amount:\n" +
	"                  $ref: \"#/components/schemas/Decimal\"\n" +
	"                exchange_rate:\n" +
	"                  $ref: \"#/components/schemas/Decimal\"\n" +
	"                currency:\n" +
	"                  $ref: \"#/components/schemas/Currency\"\n" +
	"                address:\n" +
	"                  $ref: \"#/components/schemas/Address\"\n" +
	"              required:\n" +
	"                - currency\n" +
	"                - uid\n" +
	"                - amount\n" +
	"                - exchange_rate\n" +
	"                - currency_amount\n" +
	"                - address\n" +
	"      responses:\n" +
	"        \"200\":\n" +
	"          description: OK\n" +
	"          content:\n" +
	"            application/json:\n" +
	"              schema:\n" +
	"                type: object\n" +
	"                properties:\n" +
	"                  message:\n" +
	"                    $ref: \"#/components/schemas/ResponseMessage\"\n" +
	"                  data:\n" +
	"                    type: object\n" +
	"                    properties:\n" +
	"                      withdraw_request:\n" +
	"                        type: object\n" +
	"                        properties:\n" +
	"                          id:\n" +
	"                            $ref: \"#/components/schemas/ID\"\n" +
	"                          code:\n" +
	"                            type: string\n" +
	"                            example: T-200508-170535\n" +
	"                          uid:\n" +
	"                            $ref: \"#/components/schemas/ID\"\n" +
	"                          amount:\n" +
	"                            $ref: \"#/components/schemas/Decimal\"\n" +
	"                          exchange_rate:\n" +
	"                            $ref: \"#/components/schemas/Decimal\"\n" +
	"                          currency:\n" +
	"                            $ref: \"#/components/schemas/Currency\"\n" +
	"                          currency_amount:\n" +
	"                            $ref: \"#/components/schemas/Decimal\"\n" +
	"                          address:\n" +
	"                            $ref: \"#/components/schemas/Address\"\n" +
	"                          status:\n" +
	"                            type: string\n" +
	"                          create_time:\n" +
	"                            $ref: \"#/components/schemas/Timestamp\"\n" +
	"  /v1/payment/profit/withdraw/reject/:\n" +
	"    post:\n" +
	"      summary: Reject a Profit Withdraw request\n" +
	"      tags:\n" +
	"        - Payment\n" +
	"      requestBody:\n" +
	"        required: true\n" +
	"        content:\n" +
	"          application/json:\n" +
	"            schema:\n" +
	"              type: object\n" +
	"              properties:\n" +
	"                id:\n" +
	"                  $ref: \"#/components/schemas/ID\"\n" +
	"                note:\n" +
	"                  $ref: \"#/components/schemas/Note\"\n" +
	"              required:\n" +
	"                - id\n" +
	"      responses:\n" +
	"        \"200\":\n" +
	"          description: OK\n" +
	"          content:\n" +
	"            application/json:\n" +
	"              schema:\n" +
	"                type: object\n" +
	"                properties:\n" +
	"                  message:\n" +
	"                    $ref: \"#/components/schemas/ResponseMessage\"\n" +
	"                  data:\n" +
	"                    type: object\n" +
	"                    properties:\n" +
	"                      ok:\n" +
	"                        type: boolean\n" +
	"  /v1/payment/profit/reinvest/submit/:\n" +
	"    post:\n" +
	"      summary: Add a Profit Reinvest request\n" +
	"      tags:\n" +
	"        - Payment\n" +
	"      requestBody:\n" +
	"        required: true\n" +
	"        content:\n" +
	"          application/json:\n" +
	"            schema:\n" +
	"              type: object\n" +
	"              properties:\n" +
	"                uid:\n" +
	"                  $ref: \"#/components/schemas/ID\"\n" +
	"                amount:\n" +
	"                  $ref: \"#/components/schemas/Decimal\"\n" +
	"                exchange_rate:\n" +
	"                  $ref: \"#/components/schemas/Decimal\"\n" +
	"                currency:\n" +
	"                  $ref: \"#/components/schemas/Currency\"\n" +
	"                address:\n" +
	"                  $ref: \"#/components/schemas/Address\"\n" +
	"              required:\n" +
	"                - currency\n" +
	"                - uid\n" +
	"                - amount\n" +
	"                - exchange_rate\n" +
	"                - currency_amount\n" +
	"                - address\n" +
	"      responses:\n" +
	"        \"200\":\n" +
	"          description: OK\n" +
	"          content:\n" +
	"            application/json:\n" +
	"              schema:\n" +
	"                type: object\n" +
	"                properties:\n" +
	"                  message:\n" +
	"                    $ref: \"#/components/schemas/ResponseMessage\"\n" +
	"                  data:\n" +
	"                    type: object\n" +
	"                    properties:\n" +
	"                      reinvest_request:\n" +
	"                        type: object\n" +
	"                        properties:\n" +
	"                          id:\n" +
	"                            $ref: \"#/components/schemas/ID\"\n" +
	"                          uid:\n" +
	"                            $ref: \"#/components/schemas/ID\"\n" +
	"                          amount:\n" +
	"                            $ref: \"#/components/schemas/Decimal\"\n" +
	"                          exchange_rate:\n" +
	"                            $ref: \"#/components/schemas/Decimal\"\n" +
	"                          currency:\n" +
	"                            $ref: \"#/components/schemas/Currency\"\n" +
	"                          currency_amount:\n" +
	"                            $ref: \"#/components/schemas/Decimal\"\n" +
	"                          address:\n" +
	"                            $ref: \"#/components/schemas/Address\"\n" +
	"                          status:\n" +
	"                            type: string\n" +
	"                          create_time:\n" +
	"                            $ref: \"#/components/schemas/Timestamp\"\n" +
	"  /v1/payment/profit/reinvest/approve/:\n" +
	"    post:\n" +
	"      summary: Approve a Profit Reinvest request\n" +
	"      tags:\n" +
	"        - Payment\n" +
	"      requestBody:\n" +
	"        required: true\n" +
	"        content:\n" +
	"          application/json:\n" +
	"            schema:\n" +
	"              type: object\n" +
	"              properties:\n" +
	"                id:\n" +
	"                  $ref: \"#/components/schemas/ID\"\n" +
	"                note:\n" +
	"                  $ref: \"#/components/schemas/Note\"\n" +
	"              required:\n" +
	"                - id\n" +
	"      responses:\n" +
	"        \"200\":\n" +
	"          description: OK\n" +
	"          content:\n" +
	"            application/json:\n" +
	"              schema:\n" +
	"                type: object\n" +
	"                properties:\n" +
	"                  message:\n" +
	"                    $ref: \"#/components/schemas/ResponseMessage\"\n" +
	"                  data:\n" +
	"                    type: object\n" +
	"                    properties:\n" +
	"                      ok:\n" +
	"                        type: boolean\n" +
	"  /v1/payment/profit/reinvest/reject/:\n" +
	"    post:\n" +
	"      summary: Reject a Profit Reinvest request\n" +
	"      tags:\n" +
	"        - Payment\n" +
	"      requestBody:\n" +
	"        required: true\n" +
	"        content:\n" +
	"          application/json:\n" +
	"            schema:\n" +
	"              type: object\n" +
	"              properties:\n" +
	"                id:\n" +
	"                  $ref: \"#/components/schemas/ID\"\n" +
	"                note:\n" +
	"                  $ref: \"#/components/schemas/Note\"\n" +
	"              required:\n" +
	"                - id\n" +
	"      responses:\n" +
	"        \"200\":\n" +
	"          description: OK\n" +
	"          content:\n" +
	"            application/json:\n" +
	"              schema:\n" +
	"                type: object\n" +
	"                properties:\n" +
	"                  message:\n" +
	"                    $ref: \"#/components/schemas/ResponseMessage\"\n" +
	"                  data:\n" +
	"                    type: object\n" +
	"                    properties:\n" +
	"                      ok:\n" +
	"                        type: boolean\n" +
	"  /v1/payment/p2p/transfer/:\n" +
	"    post:\n" +
	"      summary: Transfer currency to another account\n" +
	"      tags:\n" +
	"        - Payment\n" +
	"      requestBody:\n" +
	"        required: true\n" +
	"        content:\n" +
	"          application/json:\n" +
	"            schema:\n" +
	"              type: object\n" +
	"              properties:\n" +
	"                from_uid:\n" +
	"                  $ref: \"#/components/schemas/ID\"\n" +
	"                from_currency:\n" +
	"                  $ref: \"#/components/schemas/Currency\"\n" +
	"                from_amount:\n" +
	"                  $ref: \"#/components/schemas/Decimal\"\n" +
	"                to_uid:\n" +
	"                  $ref: \"#/components/schemas/ID\"\n" +
	"                to_currency:\n" +
	"                  $ref: \"#/components/schemas/Currency\"\n" +
	"                exchange_rate:\n" +
	"                  $ref: \"#/components/schemas/Decimal\"\n" +
	"                fee_amount:\n" +
	"                  $ref: \"#/components/schemas/Decimal\"\n" +
	"                note:\n" +
	"                  $ref: \"#/components/schemas/Note\"\n" +
	"              required:\n" +
	"                - from_uid\n" +
	"                - from_currency\n" +
	"                - from_amount\n" +
	"                - to_uid\n" +
	"                - to_currency\n" +
	"                - to_amount\n" +
	"                - exchange_rate\n" +
	"      responses:\n" +
	"        \"200\":\n" +
	"          description: OK\n" +
	"          content:\n" +
	"            application/json:\n" +
	"              schema:\n" +
	"                type: object\n" +
	"                properties:\n" +
	"                  message:\n" +
	"                    $ref: \"#/components/schemas/ResponseMessage\"\n" +
	"                  data:\n" +
	"                    type: object\n" +
	"                    properties:\n" +
	"                      transfer:\n" +
	"                        type: object\n" +
	"                        properties:\n" +
	"                          id:\n" +
	"                            $ref: \"#/components/schemas/ID\"\n" +
	"                          from_uid:\n" +
	"                            $ref: \"#/components/schemas/ID\"\n" +
	"                          from_currency:\n" +
	"                            $ref: \"#/components/schemas/Currency\"\n" +
	"                          from_amount:\n" +
	"                            $ref: \"#/components/schemas/Decimal\"\n" +
	"                          to_uid:\n" +
	"                            $ref: \"#/components/schemas/ID\"\n" +
	"                          to_currency:\n" +
	"                            $ref: \"#/components/schemas/Currency\"\n" +
	"                          to_amount:\n" +
	"                            $ref: \"#/components/schemas/Decimal\"\n" +
	"                          exchange_rate:\n" +
	"                            $ref: \"#/components/schemas/Decimal\"\n" +
	"                          fee_amount:\n" +
	"                            $ref: \"#/components/schemas/Decimal\"\n" +
	"                          note:\n" +
	"                            $ref: \"#/components/schemas/Note\"\n" +
	"                          status:\n" +
	"                            type: integer\n" +
	"                          create_time:\n" +
	"                            $ref: \"#/components/schemas/Timestamp\"\n" +
	"  /v1/payment/promo-code/redeem/:\n" +
	"    post:\n" +
	"      summary: Redeem a Promo code\n" +
	"      tags:\n" +
	"        - Payment\n" +
	"      requestBody:\n" +
	"        required: true\n" +
	"        content:\n" +
	"          application/json:\n" +
	"            schema:\n" +
	"              type: object\n" +
	"              properties:\n" +
	"                uid:\n" +
	"                  $ref: \"#/components/schemas/ID\"\n" +
	"                code:\n" +
	"                  type: string\n" +
	"                  example: AbcxYZ\n" +
	"              required:\n" +
	"                - uid\n" +
	"                - code\n" +
	"      responses:\n" +
	"        \"200\":\n" +
	"          description: OK\n" +
	"          content:\n" +
	"            application/json:\n" +
	"              schema:\n" +
	"                type: object\n" +
	"                properties:\n" +
	"                  message:\n" +
	"                    $ref: \"#/components/schemas/ResponseMessage\"\n" +
	"                  data:\n" +
	"                    type: object\n" +
	"                    properties:\n" +
	"                      redemption:\n" +
	"                        type: object\n" +
	"                        properties:\n" +
	"                          id:\n" +
	"                            $ref: \"#/components/schemas/ID\"\n" +
	"                          uid:\n" +
	"                            $ref: \"#/components/schemas/ID\"\n" +
	"                          code:\n" +
	"                            type: string\n" +
	"                          create_time:\n" +
	"                            $ref: \"#/components/schemas/Timestamp\"\n" +
	"  /v1/system/withdrawal/account/generate/:\n" +
	"    post:\n" +
	"      summary: Generate a new System withdrawal account\n" +
	"      tags:\n" +
	"        - System\n" +
	"      requestBody:\n" +
	"        required: true\n" +
	"        content:\n" +
	"          application/json:\n" +
	"            schema:\n" +
	"              type: object\n" +
	"              properties:\n" +
	"                currency:\n" +
	"                  $ref: \"#/components/schemas/Currency\"\n" +
	"                network:\n" +
	"                  $ref: \"#/components/schemas/Network\"\n" +
	"              required:\n" +
	"                - uid\n" +
	"      responses:\n" +
	"        \"200\":\n" +
	"          description: OK\n" +
	"          content:\n" +
	"            application/json:\n" +
	"              schema:\n" +
	"                type: object\n" +
	"                properties:\n" +
	"                  message:\n" +
	"                    $ref: \"#/components/schemas/ResponseMessage\"\n" +
	"                  data:\n" +
	"                    type: object\n" +
	"                    properties:\n" +
	"                      currency:\n" +
	"                        $ref: \"#/components/schemas/Currency\"\n" +
	"                      network:\n" +
	"                        $ref: \"#/components/schemas/Network\"\n" +
	"                      account_no:\n" +
	"                        $ref: \"#/components/schemas/Address\"\n" +
	"  /v1/system/withdrawal/account/get/:\n" +
	"    post:\n" +
	"      summary: Get the current active System withdrawal account\n" +
	"      tags:\n" +
	"        - System\n" +
	"      requestBody:\n" +
	"        required: true\n" +
	"        content:\n" +
	"          application/json:\n" +
	"            schema:\n" +
	"              type: object\n" +
	"              properties:\n" +
	"                currency:\n" +
	"                  $ref: \"#/components/schemas/Currency\"\n" +
	"                network:\n" +
	"                  $ref: \"#/components/schemas/Network\"\n" +
	"              required:\n" +
	"                - currency\n" +
	"      responses:\n" +
	"        \"200\":\n" +
	"          description: OK\n" +
	"          content:\n" +
	"            application/json:\n" +
	"              schema:\n" +
	"                type: object\n" +
	"                properties:\n" +
	"                  message:\n" +
	"                    $ref: \"#/components/schemas/ResponseMessage\"\n" +
	"                  data:\n" +
	"                    type: object\n" +
	"                    properties:\n" +
	"                      currency:\n" +
	"                        $ref: \"#/components/schemas/Currency\"\n" +
	"                      network:\n" +
	"                        $ref: \"#/components/schemas/Network\"\n" +
	"                      account_no:\n" +
	"                        $ref: \"#/components/schemas/Address\"\n" +
	"                      pull_address:\n" +
	"                        $ref: \"#/components/schemas/Address\"\n" +
	"  /v1/system/withdrawal/account/pull/:\n" +
	"    post:\n" +
	"      summary: Pull a System withdrawal address balance to a pre-defined system address\n" +
	"      tags:\n" +
	"        - System\n" +
	"      requestBody:\n" +
	"        required: true\n" +
	"        content:\n" +
	"          application/json:\n" +
	"            schema:\n" +
	"              type: object\n" +
	"              properties:\n" +
	"                currency:\n" +
	"                  $ref: \"#/components/schemas/Currency\"\n" +
	"                network:\n" +
	"                  $ref: \"#/components/schemas/Network\"\n" +
	"                account_no:\n" +
	"                  $ref: \"#/components/schemas/Address\"\n" +
	"                pull_address:\n" +
	"                  $ref: \"#/components/schemas/Address\"\n" +
	"              required:\n" +
	"                - currency\n" +
	"                - account_no\n" +
	"                - pull_address\n" +
	"      responses:\n" +
	"        \"200\":\n" +
	"          description: OK\n" +
	"          content:\n" +
	"            application/json:\n" +
	"              schema:\n" +
	"                type: object\n" +
	"                properties:\n" +
	"                  message:\n" +
	"                    $ref: \"#/components/schemas/ResponseMessage\"\n" +
	"                  data:\n" +
	"                    type: object\n" +
	"                    properties:\n" +
	"                      currency:\n" +
	"                        $ref: \"#/components/schemas/Currency\"\n" +
	"                      network:\n" +
	"                        $ref: \"#/components/schemas/Network\"\n" +
	"                      hash:\n" +
	"                        type: string\n" +
	"                        example: 947cfb43f4d3bea439f5045f9c47e811a47dae26e2976be0ef56d5e058f81e0a\n" +
	"  /v1/system/withdrawal/transfer/meta/:\n" +
	"    post:\n" +
	"      summary: Get System withdrawal meta data\n" +
	"      tags:\n" +
	"        - System\n" +
	"      responses:\n" +
	"        \"200\":\n" +
	"          description: OK\n" +
	"          content:\n" +
	"            application/json:\n" +
	"              schema:\n" +
	"                type: object\n" +
	"                properties:\n" +
	"                  message:\n" +
	"                    $ref: \"#/components/schemas/ResponseMessage\"\n" +
	"                  data:\n" +
	"                    type: object\n" +
	"                    properties:\n" +
	"                      currencies:\n" +
	"                        type: array\n" +
	"                        items:\n" +
	"                          type: object\n" +
	"                          properties:\n" +
	"                            currency:\n" +
	"                              $ref: \"#/components/schemas/Currency\"\n" +
	"                            network:\n" +
	"                              $ref: \"#/components/schemas/Network\"\n" +
	"                            fee_info:\n" +
	"                              $ref: \"#/components/schemas/FeeInfo\"\n" +
	"  /v1/system/withdrawal/transfer/submit/:\n" +
	"    post:\n" +
	"      summary: Submit a System withdrawal request\n" +
	"      tags:\n" +
	"        - System\n" +
	"      requestBody:\n" +
	"        required: true\n" +
	"        content:\n" +
	"          application/json:\n" +
	"            schema:\n" +
	"              type: object\n" +
	"              properties:\n" +
	"                request_uid:\n" +
	"                  $ref: \"#/components/schemas/ID\"\n" +
	"                currency:\n" +
	"                  $ref: \"#/components/schemas/Currency\"\n" +
	"                network:\n" +
	"                  $ref: \"#/components/schemas/Network\"\n" +
	"                src_address:\n" +
	"                  $ref: \"#/components/schemas/Address\"\n" +
	"                codes:\n" +
	"                  type: array\n" +
	"                  items:\n" +
	"                    type: string\n" +
	"                  example:\n" +
	"                    - C-200109-293319\n" +
	"                    - T-191122-540812\n" +
	"                total_amount:\n" +
	"                  $ref: \"#/components/schemas/Decimal\"\n" +
	"              required:\n" +
	"                - request_uid\n" +
	"                - currency\n" +
	"                - src_address\n" +
	"                - codes\n" +
	"                - total_amount\n" +
	"      responses:\n" +
	"        \"200\":\n" +
	"          description: OK\n" +
	"          content:\n" +
	"            application/json:\n" +
	"              schema:\n" +
	"                type: object\n" +
	"                properties:\n" +
	"                  message:\n" +
	"                    $ref: \"#/components/schemas/ResponseMessage\"\n" +
	"                  data:\n" +
	"                    type: object\n" +
	"                    properties:\n" +
	"                      request_id:\n" +
	"                        $ref: \"#/components/schemas/ID\"\n" +
	"                      address_id:\n" +
	"                        $ref: \"#/components/schemas/ID\"\n" +
	"                      status:\n" +
	"                        type: integer\n" +
	"                        example: 1\n" +
	"                      currency:\n" +
	"                        $ref: \"#/components/schemas/Currency\"\n" +
	"                      network:\n" +
	"                        $ref: \"#/components/schemas/Network\"\n" +
	"                      amount:\n" +
	"                        $ref: \"#/components/schemas/Decimal\"\n" +
	"                      amount_estimated_fee:\n" +
	"                        $ref: \"#/components/schemas/Decimal\"\n" +
	"                      combined_txn_hash:\n" +
	"                        $ref: \"#/components/schemas/TxnHash\"\n" +
	"                      create_uid:\n" +
	"                        $ref: \"#/components/schemas/ID\"\n" +
	"                      create_time:\n" +
	"                        $ref: \"#/components/schemas/Timestamp\"\n" +
	"  /v1/system/withdrawal/transfer/confirm/:\n" +
	"    post:\n" +
	"      summary: Comfirm a System withdrawal request\n" +
	"      tags:\n" +
	"        - System\n" +
	"      requestBody:\n" +
	"        required: true\n" +
	"        content:\n" +
	"          application/json:\n" +
	"            schema:\n" +
	"              type: object\n" +
	"              properties:\n" +
	"                request_id:\n" +
	"                  $ref: \"#/components/schemas/ID\"\n" +
	"              required:\n" +
	"                - request_id\n" +
	"      responses:\n" +
	"        \"200\":\n" +
	"          description: OK\n" +
	"          content:\n" +
	"            application/json:\n" +
	"              schema:\n" +
	"                type: object\n" +
	"                properties:\n" +
	"                  message:\n" +
	"                    $ref: \"#/components/schemas/ResponseMessage\"\n" +
	"  /v1/system/withdrawal/transfer/replace/:\n" +
	"    post:\n" +
	"      summary: Replace a System withdrawal request\n" +
	"      tags:\n" +
	"        - System\n" +
	"      requestBody:\n" +
	"        required: true\n" +
	"        content:\n" +
	"          application/json:\n" +
	"            schema:\n" +
	"              type: object\n" +
	"              properties:\n" +
	"                request_id:\n" +
	"                  $ref: \"#/components/schemas/ID\"\n" +
	"                fee_info:\n" +
	"                  $ref: \"#/components/schemas/FeeInfo\"\n" +
	"              required:\n" +
	"                - request_id\n" +
	"                - fee_info\n" +
	"      responses:\n" +
	"        \"200\":\n" +
	"          description: OK\n" +
	"          content:\n" +
	"            application/json:\n" +
	"              schema:\n" +
	"                type: object\n" +
	"                properties:\n" +
	"                  message:\n" +
	"                    $ref: \"#/components/schemas/ResponseMessage\"\n" +
	"  /v1/system/pool/bonus/leader/checkout/:\n" +
	"    post:\n" +
	"      summary: Checkout Leader Bonus Pool meta data\n" +
	"      tags:\n" +
	"        - System\n" +
	"        - Bonus Pool\n" +
	"      requestBody:\n" +
	"        required: true\n" +
	"        content:\n" +
	"          application/json:\n" +
	"            schema:\n" +
	"              type: object\n" +
	"              properties:\n" +
	"                from_date:\n" +
	"                  $ref: \"#/components/schemas/Date\"\n" +
	"                to_date:\n" +
	"                  $ref: \"#/components/schemas/Date\"\n" +
	"              required:\n" +
	"                - from_date\n" +
	"                - to_date\n" +
	"      responses:\n" +
	"        \"200\":\n" +
	"          description: OK\n" +
	"          content:\n" +
	"            application/json:\n" +
	"              schema:\n" +
	"                type: object\n" +
	"                properties:\n" +
	"                  message:\n" +
	"                    $ref: \"#/components/schemas/ResponseMessage\"\n" +
	"                  data:\n" +
	"                    type: object\n" +
	"                    properties:\n" +
	"                      execution_hash:\n" +
	"                        $ref: \"#/components/schemas/TxnHash\"\n" +
	"                      from_date:\n" +
	"                        $ref: \"#/components/schemas/Date\"\n" +
	"                      to_date:\n" +
	"                        $ref: \"#/components/schemas/Date\"\n" +
	"                      total_amount:\n" +
	"                        $ref: \"#/components/schemas/Decimal\"\n" +
	"                      tier_info_list:\n" +
	"                        type: array\n" +
	"                        items:\n" +
	"                          $ref: \"#/components/schemas/BonusPoolLeaderTierInfo\"\n" +
	"                    required:\n" +
	"                      - execution_hash\n" +
	"                      - from_date\n" +
	"                      - to_date\n" +
	"                      - total_amount\n" +
	"                      - tier_info_list\n" +
	"  /v1/system/pool/bonus/leader/execute/:\n" +
	"    post:\n" +
	"      summary: Execute a Leader Bonus Pool reward session\n" +
	"      tags:\n" +
	"        - System\n" +
	"        - Bonus Pool\n" +
	"      requestBody:\n" +
	"        required: true\n" +
	"        content:\n" +
	"          application/json:\n" +
	"            schema:\n" +
	"              type: object\n" +
	"              properties:\n" +
	"                execution_hash:\n" +
	"                  $ref: \"#/components/schemas/TxnHash\"\n" +
	"                from_date:\n" +
	"                  $ref: \"#/components/schemas/Date\"\n" +
	"                to_date:\n" +
	"                  $ref: \"#/components/schemas/Date\"\n" +
	"                total_amount:\n" +
	"                  $ref: \"#/components/schemas/Decimal\"\n" +
	"                tier_info_list:\n" +
	"                  type: array\n" +
	"                  items:\n" +
	"                    $ref: \"#/components/schemas/BonusPoolLeaderTierInfo\"\n" +
	"              required:\n" +
	"                - execution_hash\n" +
	"                - from_date\n" +
	"                - to_date\n" +
	"                - total_amount\n" +
	"                - tier_info_list\n" +
	"      responses:\n" +
	"        \"200\":\n" +
	"          description: OK\n" +
	"          content:\n" +
	"            application/json:\n" +
	"              schema:\n" +
	"                type: object\n" +
	"                properties:\n" +
	"                  message:\n" +
	"                    $ref: \"#/components/schemas/ResponseMessage\"\n" +
	"                  data:\n" +
	"                    type: object\n" +
	"                    properties:\n" +
	"                      execution:\n" +
	"                        type: object\n" +
	"                        properties:\n" +
	"                          id:\n" +
	"                            $ref: \"#/components/schemas/ID\"\n" +
	"                          from_date:\n" +
	"                            $ref: \"#/components/schemas/Date\"\n" +
	"                          to_date:\n" +
	"                            $ref: \"#/components/schemas/Date\"\n" +
	"                          total_amount:\n" +
	"                            $ref: \"#/components/schemas/Decimal\"\n" +
	"                          status:\n" +
	"                            $ref: \"#/components/schemas/BonusPoolLeaderExecutionStatus\"\n" +
	"                          create_time:\n" +
	"                            $ref: \"#/components/schemas/Timestamp\"\n" +
	"                          update_time:\n" +
	"                            $ref: \"#/components/schemas/Timestamp\"\n" +
	"                        required:\n" +
	"                          - id\n" +
	"                          - from_date\n" +
	"                          - to_date\n" +
	"                          - total_amount\n" +
	"                          - status\n" +
	"                          - create_time\n" +
	"                          - update_time\n" +
	""
