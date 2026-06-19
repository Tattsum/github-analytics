/* eslint-disable */
import type { TypedDocumentNode as DocumentNode } from '@graphql-typed-document-node/core';
export type Maybe<T> = T | null;
export type InputMaybe<T> = Maybe<T>;
export type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]?: Maybe<T[SubKey]> };
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]: Maybe<T[SubKey]> };
export type MakeEmpty<T extends { [key: string]: unknown }, K extends keyof T> = { [_ in K]?: never };
export type Incremental<T> = T | { [P in keyof T]?: P extends ' $fragmentName' | '__typename' ? T[P] : never };
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: { input: string; output: string; }
  String: { input: string; output: string; }
  Boolean: { input: boolean; output: boolean; }
  Int: { input: number; output: number; }
  Float: { input: number; output: number; }
};

export type MemberStats = {
  __typename?: 'MemberStats';
  login: Scalars['String']['output'];
  name: Scalars['String']['output'];
  prToReviewRatio: Scalars['Float']['output'];
  totalAdditions: Scalars['Int']['output'];
  totalCommits: Scalars['Int']['output'];
  totalDeletions: Scalars['Int']['output'];
  totalIssues: Scalars['Int']['output'];
  totalPRCreated: Scalars['Int']['output'];
  totalPRMerged: Scalars['Int']['output'];
  totalReviews: Scalars['Int']['output'];
};

export type Query = {
  __typename?: 'Query';
  member?: Maybe<UserStatistics>;
  members: Array<MemberStats>;
  repositories: Array<RepositoryStats>;
  repository?: Maybe<RepositoryStats>;
  teamSummary: TeamSummary;
};


export type QueryMemberArgs = {
  login: Scalars['String']['input'];
};


export type QueryRepositoryArgs = {
  nameWithOwner: Scalars['String']['input'];
};

export type RepositoryActivity = {
  __typename?: 'RepositoryActivity';
  commitCount: Scalars['Int']['output'];
  firstActivity: Scalars['String']['output'];
  issueCount: Scalars['Int']['output'];
  lastActivity: Scalars['String']['output'];
  prCount: Scalars['Int']['output'];
  repository: Scalars['String']['output'];
  reviewCount: Scalars['Int']['output'];
  totalAdditions: Scalars['Int']['output'];
  totalDeletions: Scalars['Int']['output'];
};

export type RepositoryContributor = {
  __typename?: 'RepositoryContributor';
  additions: Scalars['Int']['output'];
  commitCount: Scalars['Int']['output'];
  deletions: Scalars['Int']['output'];
  login: Scalars['String']['output'];
  prCreated: Scalars['Int']['output'];
  reviewCount: Scalars['Int']['output'];
};

export type RepositoryStats = {
  __typename?: 'RepositoryStats';
  contributorCount: Scalars['Int']['output'];
  contributors: Array<RepositoryContributor>;
  nameWithOwner: Scalars['String']['output'];
  total: RepositoryTotals;
};

export type RepositoryTotals = {
  __typename?: 'RepositoryTotals';
  additions: Scalars['Int']['output'];
  commits: Scalars['Int']['output'];
  deletions: Scalars['Int']['output'];
  issues: Scalars['Int']['output'];
  prCreated: Scalars['Int']['output'];
  prMerged: Scalars['Int']['output'];
  reviews: Scalars['Int']['output'];
};

export type RoleTransitionPoint = {
  __typename?: 'RoleTransitionPoint';
  description: Scalars['String']['output'];
  prCreated: Scalars['Int']['output'];
  ratio: Scalars['Float']['output'];
  reviewCount: Scalars['Int']['output'];
  year: Scalars['Int']['output'];
};

export type TeamSummary = {
  __typename?: 'TeamSummary';
  memberCount: Scalars['Int']['output'];
  repositoryCount: Scalars['Int']['output'];
  totalAdditions: Scalars['Int']['output'];
  totalCommits: Scalars['Int']['output'];
  totalDeletions: Scalars['Int']['output'];
  totalIssues: Scalars['Int']['output'];
  totalPRCreated: Scalars['Int']['output'];
  totalPRMerged: Scalars['Int']['output'];
  totalReviews: Scalars['Int']['output'];
};

export type UserStatistics = {
  __typename?: 'UserStatistics';
  firstActivityYear: Scalars['Int']['output'];
  login: Scalars['String']['output'];
  longTermRepositories: Array<RepositoryActivity>;
  name: Scalars['String']['output'];
  peakActivityCommits: Scalars['Int']['output'];
  peakActivityYear: Scalars['Int']['output'];
  prToReviewRatio: Scalars['Float']['output'];
  roleTransition: Array<RoleTransitionPoint>;
  topRepositories: Array<RepositoryActivity>;
  totalAdditions: Scalars['Int']['output'];
  totalCommits: Scalars['Int']['output'];
  totalDeletions: Scalars['Int']['output'];
  totalIssues: Scalars['Int']['output'];
  totalPRCreated: Scalars['Int']['output'];
  totalPRMerged: Scalars['Int']['output'];
  totalReviews: Scalars['Int']['output'];
  yearlyStats: Array<YearlyStatistics>;
};

export type YearlyStatistics = {
  __typename?: 'YearlyStatistics';
  commitCount: Scalars['Int']['output'];
  issueCount: Scalars['Int']['output'];
  prCreated: Scalars['Int']['output'];
  prMerged: Scalars['Int']['output'];
  reviewCount: Scalars['Int']['output'];
  totalAdditions: Scalars['Int']['output'];
  totalDeletions: Scalars['Int']['output'];
  year: Scalars['Int']['output'];
};

export type MemberDetailQueryVariables = Exact<{
  login: Scalars['String']['input'];
}>;


export type MemberDetailQuery = { __typename?: 'Query', member?: { __typename?: 'UserStatistics', login: string, name: string, totalCommits: number, totalPRCreated: number, totalPRMerged: number, totalReviews: number, totalIssues: number, totalAdditions: number, totalDeletions: number, prToReviewRatio: number, firstActivityYear: number, peakActivityYear: number, peakActivityCommits: number, yearlyStats: Array<{ __typename?: 'YearlyStatistics', year: number, commitCount: number, prCreated: number, prMerged: number, reviewCount: number, issueCount: number, totalAdditions: number, totalDeletions: number }>, topRepositories: Array<{ __typename?: 'RepositoryActivity', repository: string, commitCount: number, prCount: number, reviewCount: number, issueCount: number, totalAdditions: number, totalDeletions: number, firstActivity: string, lastActivity: string }>, longTermRepositories: Array<{ __typename?: 'RepositoryActivity', repository: string, commitCount: number, prCount: number, reviewCount: number, issueCount: number, totalAdditions: number, totalDeletions: number, firstActivity: string, lastActivity: string }>, roleTransition: Array<{ __typename?: 'RoleTransitionPoint', year: number, description: string, prCreated: number, reviewCount: number, ratio: number }> } | null };

export type RepositoriesQueryVariables = Exact<{ [key: string]: never; }>;


export type RepositoriesQuery = { __typename?: 'Query', repositories: Array<{ __typename?: 'RepositoryStats', nameWithOwner: string, contributorCount: number, total: { __typename?: 'RepositoryTotals', commits: number, prCreated: number, prMerged: number, issues: number, reviews: number, additions: number, deletions: number } }> };

export type RepositoryQueryVariables = Exact<{
  nameWithOwner: Scalars['String']['input'];
}>;


export type RepositoryQuery = { __typename?: 'Query', repository?: { __typename?: 'RepositoryStats', nameWithOwner: string, contributorCount: number, total: { __typename?: 'RepositoryTotals', commits: number, prCreated: number, prMerged: number, issues: number, reviews: number, additions: number, deletions: number }, contributors: Array<{ __typename?: 'RepositoryContributor', login: string, commitCount: number, prCreated: number, reviewCount: number, additions: number, deletions: number }> } | null };

export type TeamOverviewSummaryQueryVariables = Exact<{ [key: string]: never; }>;


export type TeamOverviewSummaryQuery = { __typename?: 'Query', teamSummary: { __typename?: 'TeamSummary', memberCount: number, repositoryCount: number, totalCommits: number, totalPRCreated: number, totalPRMerged: number, totalIssues: number, totalReviews: number, totalAdditions: number, totalDeletions: number } };

export type TeamOverviewMembersQueryVariables = Exact<{ [key: string]: never; }>;


export type TeamOverviewMembersQuery = { __typename?: 'Query', members: Array<{ __typename?: 'MemberStats', login: string, name: string, totalCommits: number, totalPRCreated: number, totalPRMerged: number, totalIssues: number, totalReviews: number, totalAdditions: number, totalDeletions: number, prToReviewRatio: number }> };


export const MemberDetailDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"MemberDetail"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"login"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"String"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"member"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"login"},"value":{"kind":"Variable","name":{"kind":"Name","value":"login"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"login"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"totalCommits"}},{"kind":"Field","name":{"kind":"Name","value":"totalPRCreated"}},{"kind":"Field","name":{"kind":"Name","value":"totalPRMerged"}},{"kind":"Field","name":{"kind":"Name","value":"totalReviews"}},{"kind":"Field","name":{"kind":"Name","value":"totalIssues"}},{"kind":"Field","name":{"kind":"Name","value":"totalAdditions"}},{"kind":"Field","name":{"kind":"Name","value":"totalDeletions"}},{"kind":"Field","name":{"kind":"Name","value":"prToReviewRatio"}},{"kind":"Field","name":{"kind":"Name","value":"firstActivityYear"}},{"kind":"Field","name":{"kind":"Name","value":"peakActivityYear"}},{"kind":"Field","name":{"kind":"Name","value":"peakActivityCommits"}},{"kind":"Field","name":{"kind":"Name","value":"yearlyStats"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"year"}},{"kind":"Field","name":{"kind":"Name","value":"commitCount"}},{"kind":"Field","name":{"kind":"Name","value":"prCreated"}},{"kind":"Field","name":{"kind":"Name","value":"prMerged"}},{"kind":"Field","name":{"kind":"Name","value":"reviewCount"}},{"kind":"Field","name":{"kind":"Name","value":"issueCount"}},{"kind":"Field","name":{"kind":"Name","value":"totalAdditions"}},{"kind":"Field","name":{"kind":"Name","value":"totalDeletions"}}]}},{"kind":"Field","name":{"kind":"Name","value":"topRepositories"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"repository"}},{"kind":"Field","name":{"kind":"Name","value":"commitCount"}},{"kind":"Field","name":{"kind":"Name","value":"prCount"}},{"kind":"Field","name":{"kind":"Name","value":"reviewCount"}},{"kind":"Field","name":{"kind":"Name","value":"issueCount"}},{"kind":"Field","name":{"kind":"Name","value":"totalAdditions"}},{"kind":"Field","name":{"kind":"Name","value":"totalDeletions"}},{"kind":"Field","name":{"kind":"Name","value":"firstActivity"}},{"kind":"Field","name":{"kind":"Name","value":"lastActivity"}}]}},{"kind":"Field","name":{"kind":"Name","value":"longTermRepositories"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"repository"}},{"kind":"Field","name":{"kind":"Name","value":"commitCount"}},{"kind":"Field","name":{"kind":"Name","value":"prCount"}},{"kind":"Field","name":{"kind":"Name","value":"reviewCount"}},{"kind":"Field","name":{"kind":"Name","value":"issueCount"}},{"kind":"Field","name":{"kind":"Name","value":"totalAdditions"}},{"kind":"Field","name":{"kind":"Name","value":"totalDeletions"}},{"kind":"Field","name":{"kind":"Name","value":"firstActivity"}},{"kind":"Field","name":{"kind":"Name","value":"lastActivity"}}]}},{"kind":"Field","name":{"kind":"Name","value":"roleTransition"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"year"}},{"kind":"Field","name":{"kind":"Name","value":"description"}},{"kind":"Field","name":{"kind":"Name","value":"prCreated"}},{"kind":"Field","name":{"kind":"Name","value":"reviewCount"}},{"kind":"Field","name":{"kind":"Name","value":"ratio"}}]}}]}}]}}]} as unknown as DocumentNode<MemberDetailQuery, MemberDetailQueryVariables>;
export const RepositoriesDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"Repositories"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"repositories"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"nameWithOwner"}},{"kind":"Field","name":{"kind":"Name","value":"contributorCount"}},{"kind":"Field","name":{"kind":"Name","value":"total"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"commits"}},{"kind":"Field","name":{"kind":"Name","value":"prCreated"}},{"kind":"Field","name":{"kind":"Name","value":"prMerged"}},{"kind":"Field","name":{"kind":"Name","value":"issues"}},{"kind":"Field","name":{"kind":"Name","value":"reviews"}},{"kind":"Field","name":{"kind":"Name","value":"additions"}},{"kind":"Field","name":{"kind":"Name","value":"deletions"}}]}}]}}]}}]} as unknown as DocumentNode<RepositoriesQuery, RepositoriesQueryVariables>;
export const RepositoryDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"Repository"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"nameWithOwner"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"String"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"repository"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"nameWithOwner"},"value":{"kind":"Variable","name":{"kind":"Name","value":"nameWithOwner"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"nameWithOwner"}},{"kind":"Field","name":{"kind":"Name","value":"contributorCount"}},{"kind":"Field","name":{"kind":"Name","value":"total"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"commits"}},{"kind":"Field","name":{"kind":"Name","value":"prCreated"}},{"kind":"Field","name":{"kind":"Name","value":"prMerged"}},{"kind":"Field","name":{"kind":"Name","value":"issues"}},{"kind":"Field","name":{"kind":"Name","value":"reviews"}},{"kind":"Field","name":{"kind":"Name","value":"additions"}},{"kind":"Field","name":{"kind":"Name","value":"deletions"}}]}},{"kind":"Field","name":{"kind":"Name","value":"contributors"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"login"}},{"kind":"Field","name":{"kind":"Name","value":"commitCount"}},{"kind":"Field","name":{"kind":"Name","value":"prCreated"}},{"kind":"Field","name":{"kind":"Name","value":"reviewCount"}},{"kind":"Field","name":{"kind":"Name","value":"additions"}},{"kind":"Field","name":{"kind":"Name","value":"deletions"}}]}}]}}]}}]} as unknown as DocumentNode<RepositoryQuery, RepositoryQueryVariables>;
export const TeamOverviewSummaryDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"TeamOverviewSummary"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"teamSummary"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"memberCount"}},{"kind":"Field","name":{"kind":"Name","value":"repositoryCount"}},{"kind":"Field","name":{"kind":"Name","value":"totalCommits"}},{"kind":"Field","name":{"kind":"Name","value":"totalPRCreated"}},{"kind":"Field","name":{"kind":"Name","value":"totalPRMerged"}},{"kind":"Field","name":{"kind":"Name","value":"totalIssues"}},{"kind":"Field","name":{"kind":"Name","value":"totalReviews"}},{"kind":"Field","name":{"kind":"Name","value":"totalAdditions"}},{"kind":"Field","name":{"kind":"Name","value":"totalDeletions"}}]}}]}}]} as unknown as DocumentNode<TeamOverviewSummaryQuery, TeamOverviewSummaryQueryVariables>;
export const TeamOverviewMembersDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"TeamOverviewMembers"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"members"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"login"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"totalCommits"}},{"kind":"Field","name":{"kind":"Name","value":"totalPRCreated"}},{"kind":"Field","name":{"kind":"Name","value":"totalPRMerged"}},{"kind":"Field","name":{"kind":"Name","value":"totalIssues"}},{"kind":"Field","name":{"kind":"Name","value":"totalReviews"}},{"kind":"Field","name":{"kind":"Name","value":"totalAdditions"}},{"kind":"Field","name":{"kind":"Name","value":"totalDeletions"}},{"kind":"Field","name":{"kind":"Name","value":"prToReviewRatio"}}]}}]}}]} as unknown as DocumentNode<TeamOverviewMembersQuery, TeamOverviewMembersQueryVariables>;