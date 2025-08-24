# pull model
curl http://localhost:11434/api/pull -d '{"model":"gpt-oss:20b}'

# run model
ollama run gpt-oss:20b

# show loaded models
ollama ps

# stop loaded model
ollama stop gpt-oss:20b

# generate text
curl http://localhost:11434/api/generate -d '{"model":"gpt-oss:20b", "prompt":"請推論 OpenAI 為什麼會推出 gpt-oss 這樣的開源權重模型？背後有何策略？"}'

# api key: 0c43a443c9d446cea83ed5220eb2c65d.ZAit13hXv3uIxKJTmkYi8mLB

# 介紹的非常棒
# https://www.ernestchiang.com/zh/notes/ai/openai-gpt-oss/

# llm speed check
# https://www.llmspeedcheck.com/

# human last exam
# https://zhuanlan.zhihu.com/p/21151293555
# leader board: https://artificialanalysis.ai/evaluations/humanitys-last-exam
# https://huggingface.co/datasets/cais/hle
# Consider the language $L$ defined by the regular expression $( (b | c)^* a ( a | b c | c b | b b b )^* (b a | b b a | c a) )^* (b | c)^* a ( a | b c | c b | b b b )^*$. How many states does the minimal deterministic finite-state automaton that recognizes this language have?

# AIME
# 

# Answer Choices:
# A. 1
# B. 2
# C. 3
# D. 4
# E. 5
# F. 6
# G. 7
# H. 8
# I. 9
# J. 10