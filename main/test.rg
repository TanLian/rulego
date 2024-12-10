fn backtracking(nums, sum, start, result, target, candidates) {
    if sum == target {
        result.Push(nums.Clone());
        return 1;
    }

    if sum > target {
        return 1;
    }

    for i = start; i < candidates.Len(); i++ {
        nums.Push(candidates[i]);
        sum += candidates[i];
        backtracking(nums, sum, i, result, target, candidates);
        sum -= candidates[i];
        nums.Pop();
    }
}

fn combinationSum(candidates, target) {
    result = [];
    backtracking([], 0, 0, result, target, candidates);
    return result;
}

println(combinationSum([2,3,7], 7));