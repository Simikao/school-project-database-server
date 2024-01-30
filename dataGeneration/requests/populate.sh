baseurl="localhost:8080"
current_path=$(pwd)
# Make a post
function new_user() {
    for x in $(seq 0 19); do
        jq ".[$x]" < $current_path/../users.json \
            | curl -X POST "$baseurl/u" \
            -H 'Content-Type: application/json' \
            --data-binary '@-' \
            | jq
    done

    # curl -X POST "$baseurl/add"            \
    #      -H 'Content-Type: application/json' \
    #      -d '@../users.json'                \
    #      | jq
}

function new_community() {
    for x in $(seq 0 12); do 
    jq ".[$x]" < $current_path/../community.json \
        | curl -X POST "$baseurl/new-community" \
        -H 'Content-Type: application/json' \
        --data-binary '@-' \
        | jq
    done
}

function new_post() {
    for x in $(seq 0 12); do
    jq ".[$x]" < $current_path/../posts.json \
        | curl -X POST "$baseurl/new-post" \
        -H 'Content-Type: application/json' \
        --data-binary '@-' \
        | jq
    done
}


function add_all() {
    new_user
    sleep .5
    new_post
    sleep .5
    new_community
}
# new_user