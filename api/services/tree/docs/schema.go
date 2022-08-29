package docs

var doc = "" +
	"openapi: \"3.0.0\"\n" +
	"info:\n" +
	"  version: 1.0.0\n" +
	"  title: Torque Tree Service\n" +
	"servers:\n" +
	"  - url: \"/\"\n" +
	"components:\n" +
	"  schemas:\n" +
	"    ID:\n" +
	"      type: integer\n" +
	"      minimum: 1\n" +
	"      example: 8466032601\n" +
	"    Decimal:\n" +
	"      type: string\n" +
	"      description: Decimal number in string\n" +
	"      example: \"14211.364041\"\n" +
	"    Errors:\n" +
	"      type: array\n" +
	"      items:\n" +
	"        type: string\n" +
	"        example: Field validation for 'user_id' failed on the 'required' tag.\n" +
	"    Node:\n" +
	"      type: object\n" +
	"      properties:\n" +
	"        user:\n" +
	"          type: object\n" +
	"          properties:\n" +
	"            id:\n" +
	"              $ref: \"#/components/schemas/ID\"\n" +
	"            username:\n" +
	"              type: string\n" +
	"              example: helloworld\n" +
	"            first_name:\n" +
	"              type: string\n" +
	"              example: world\n" +
	"            last_name:\n" +
	"              type: string\n" +
	"              example: hello\n" +
	"            email:\n" +
	"              type: string\n" +
	"              example: abc@gmail.com\n" +
	"            referral_code:\n" +
	"              type: integer\n" +
	"              example: 987654321\n" +
	"            type:\n" +
	"              type: integer\n" +
	"              enum:\n" +
	"                - 1\n" +
	"                - 2\n" +
	"                - 3\n" +
	"                - 4\n" +
	"                - 5\n" +
	"                - 6\n" +
	"            create_time:\n" +
	"              type: integer\n" +
	"              example: 1581354000\n" +
	"          required:\n" +
	"            - id\n" +
	"            - username\n" +
	"            - first_name\n" +
	"            - last_name\n" +
	"            - email\n" +
	"            - referral_code\n" +
	"            - type\n" +
	"            - create_time\n" +
	"        user_id:\n" +
	"          type: integer\n" +
	"        username:\n" +
	"          type: string\n" +
	"          example: helloworld\n" +
	"        user_role:\n" +
	"          type: string\n" +
	"          example: agent\n" +
	"        balance_usd:\n" +
	"          $ref: \"#/components/schemas/Decimal\"\n" +
	"        coin_balance_map:\n" +
	"          $ref: \"#/components/schemas/CoinBalanceMap\"\n" +
	"        children_balance_usd:\n" +
	"          $ref: \"#/components/schemas/Decimal\"\n" +
	"        coin_children_balance_map:\n" +
	"          $ref: \"#/components/schemas/CoinBalanceMap\"\n" +
	"        descendants_count:\n" +
	"          type: integer\n" +
	"        level:\n" +
	"          type: integer\n" +
	"        children:\n" +
	"          type: array\n" +
	"          items:\n" +
	"            type: object\n" +
	"            nullable: true\n" +
	"            $ref: \"#/components/schemas/Node\"\n" +
	"      required:\n" +
	"        - user\n" +
	"        - user_id\n" +
	"        - username\n" +
	"        - user_role\n" +
	"        - level\n" +
	"    CoinBalanceMap:\n" +
	"      type: object\n" +
	"      properties:\n" +
	"        LTC:\n" +
	"          $ref: \"#/components/schemas/Decimal\"\n" +
	"        BTC:\n" +
	"          $ref: \"#/components/schemas/Decimal\"\n" +
	"        ETH:\n" +
	"          $ref: \"#/components/schemas/Decimal\"\n" +
	"        USDT:\n" +
	"          $ref: \"#/components/schemas/Decimal\"\n" +
	"      required:\n" +
	"        - LTC\n" +
	"        - BTC\n" +
	"        - ETH\n" +
	"        - USDT\n" +
	"paths:\n" +
	"  /v1/basic/get_node_down/:\n" +
	"    post:\n" +
	"      summary: Collect info on user's down tree\n" +
	"      tags:\n" +
	"        - Basic\n" +
	"      requestBody:\n" +
	"        required: true\n" +
	"        content:\n" +
	"          application/json:\n" +
	"            schema:\n" +
	"              type: object\n" +
	"              properties:\n" +
	"                user_id:\n" +
	"                  $ref: \"#/components/schemas/ID\"\n" +
	"                options:\n" +
	"                  type: object\n" +
	"                  properties:\n" +
	"                    root_uid:\n" +
	"                      $ref: \"#/components/schemas/ID\"\n" +
	"                    limit_level:\n" +
	"                      type: integer\n" +
	"                      description: Depth limit of the result tree\n" +
	"                      example: 7\n" +
	"                    get_coin_map:\n" +
	"                      type: boolean\n" +
	"                    get_children_coin_map:\n" +
	"                      type: boolean\n" +
	"                    use_raw_coin_map:\n" +
	"                      type: boolean\n" +
	"                    fetch_root_only:\n" +
	"                      type: boolean\n" +
	"              required:\n" +
	"                - user_id\n" +
	"      responses:\n" +
	"        \"200\":\n" +
	"          description: OK\n" +
	"          content:\n" +
	"            application/json:\n" +
	"              schema:\n" +
	"                type: object\n" +
	"                properties:\n" +
	"                  errors:\n" +
	"                    $ref: \"#/components/schemas/Errors\"\n" +
	"                  data:\n" +
	"                    type: object\n" +
	"                    nullable: true\n" +
	"                    properties:\n" +
	"                      node:\n" +
	"                        $ref: \"#/components/schemas/Node\"\n" +
	"                    required:\n" +
	"                      - node\n" +
	"  /v1/promotion/stats/get/:\n" +
	"    post:\n" +
	"      summary: Get user promotion stats\n" +
	"      tags:\n" +
	"        - Promotion\n" +
	"      requestBody:\n" +
	"        required: true\n" +
	"        content:\n" +
	"          application/json:\n" +
	"            schema:\n" +
	"              type: object\n" +
	"              properties:\n" +
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
	"                  errors:\n" +
	"                    $ref: \"#/components/schemas/Errors\"\n" +
	"                  data:\n" +
	"                    type: object\n" +
	"                    nullable: true\n" +
	"                    properties:\n" +
	"                      stats:\n" +
	"                        type: object\n" +
	"                        properties:\n" +
	"                          uid:\n" +
	"                            $ref: \"#/components/schemas/ID\"\n" +
	"                          down_line_stats_list:\n" +
	"                            type: array\n" +
	"                            items:\n" +
	"                              type: object\n" +
	"                              properties:\n" +
	"                                from_tier:\n" +
	"                                  type: string\n" +
	"                                  example: global_partner\n" +
	"                                to_tier:\n" +
	"                                  type: string\n" +
	"                                  example: mentor\n" +
	"                                uids:\n" +
	"                                  type: array\n" +
	"                                  items:\n" +
	"                                    $ref: \"#/components/schemas/ID\"\n" +
	"                              required:\n" +
	"                                - from_tier\n" +
	"                                - to_tier\n" +
	"                                - uids\n" +
	"                        required:\n" +
	"                          - uid\n" +
	"                          - down_line_stats_list\n" +
	"                    required:\n" +
	"                      - stats\n" +
	""
